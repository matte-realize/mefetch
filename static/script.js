let uploadedImage = null;
let importedAscii = null;
let dragSrc = null;
let previewURL = null;
let previewSeq = 0;
let previewTimer = null;

window.onload = () => {
  addSection('Info');
  addField('OS', '');
  addField('Shell', '');
  addField('Editor', '');
  updatePreview();
};

function getFields() {
  const rows = document.querySelectorAll('#fields-list > div');
  const fields = [];
  rows.forEach(row => {
    if (row.dataset.type === 'spacer') {
      fields.push({ isSpacer: true });
      return;
    }
    const labelEl = row.querySelector('.field-label');
    const valueEl = row.querySelector('.field-value');
    if (!labelEl) return;
    const label = labelEl.value.trim();
    const value = valueEl ? valueEl.value.trim() : '';
    if (!label) return;
    if (value === '__divider__') {
      fields.push({ label, isDivider: true });
    } else if (value) {
      fields.push({ label, value, isDivider: false });
    }
  });
  return fields;
}

function buildForm() {
  const form = new FormData();
  if (uploadedImage) {
    form.append('image', uploadedImage);
  } else if (importedAscii) {
    form.append('ascii', importedAscii);
  }
  form.append('username',   document.getElementById('username').value.trim());
  form.append('hostname',   document.getElementById('hostname').value.trim());
  form.append('background',  document.getElementById('background').value);
  form.append('keycolor',    document.getElementById('keycolor').value);
  form.append('textcolor',   document.getElementById('textcolor').value);
  form.append('showstats',   document.getElementById('showstats').checked);

  getFields().forEach(f => {
    if (f.isSpacer) {
      form.append('field', '---space---');
    } else if (f.isDivider) {
      form.append('field', `---:${f.label}`);
    } else {
      form.append('field', `${f.label}:${f.value}`);
    }
  });

  return form;
}

function updatePreview() {
  clearTimeout(previewTimer);
  previewTimer = setTimeout(renderPreview, 250);
}

async function renderPreview() {
  const seq = ++previewSeq;
  try {
    const resp = await fetch('/card/generate', { method: 'POST', body: buildForm() });
    if (!resp.ok) return;
    const blob = await resp.blob();
    if (seq !== previewSeq) return;
    if (previewURL) URL.revokeObjectURL(previewURL);
    previewURL = URL.createObjectURL(blob);
    document.getElementById('card-preview').src = previewURL;
  } catch (e) {
  }
}

function addField(label = '', value = '') {
  const list = document.getElementById('fields-list');
  const row = document.createElement('div');
  row.className = 'field-row';
  row.draggable = true;
  row.dataset.type = 'field';
  row.innerHTML = `
    <span class="drag-handle" title="drag to reorder">⠿</span>
    <input type="text" class="field-label" placeholder="label" oninput="updatePreview()">
    <input type="text" class="field-value" placeholder="value" oninput="updatePreview()">
    <button class="btn danger" onclick="removeField(this)">✕</button>
  `;
  row.querySelector('.field-label').value = label;
  row.querySelector('.field-value').value = value;
  addDragListeners(row);
  list.appendChild(row);
  updatePreview();
}

function addSection(label = '') {
  const list = document.getElementById('fields-list');
  const row = document.createElement('div');
  row.className = 'field-row';
  row.draggable = true;
  row.dataset.type = 'section';
  row.innerHTML = `
    <span class="drag-handle" title="drag to reorder">⠿</span>
    <input type="text" class="field-label" placeholder="section name"
           oninput="updatePreview()" style="grid-column: span 2;">
    <input type="hidden" class="field-value" value="__divider__">
    <button class="btn danger" onclick="removeField(this)">✕</button>
  `;
  row.querySelector('.field-label').value = label;
  addDragListeners(row);
  list.appendChild(row);
  updatePreview();
}

function addSpacer() {
  const list = document.getElementById('fields-list');
  const row = document.createElement('div');
  row.className = 'spacer-row';
  row.draggable = true;
  row.dataset.type = 'spacer';
  row.innerHTML = `
    <span class="drag-handle" title="drag to reorder">⠿</span>
    <div class="spacer-label">— empty space —</div>
    <button class="btn danger" onclick="removeField(this)">✕</button>
  `;
  addDragListeners(row);
  list.appendChild(row);
  updatePreview();
}

function removeField(btn) {
  btn.closest('div').remove();
  updatePreview();
}

function addDragListeners(row) {
  row.addEventListener('dragstart', e => {
    dragSrc = row;
    row.classList.add('dragging');
    e.dataTransfer.effectAllowed = 'move';
  });

  row.addEventListener('dragend', () => {
    row.classList.remove('dragging');
    document.querySelectorAll('.drag-over').forEach(el => {
      el.classList.remove('drag-over');
    });
    updatePreview();
  });

  row.addEventListener('dragover', e => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
    if (row !== dragSrc) {
      row.classList.add('drag-over');
    }
  });

  row.addEventListener('dragleave', () => {
    row.classList.remove('drag-over');
  });

  row.addEventListener('drop', e => {
    e.preventDefault();
    row.classList.remove('drag-over');
    if (dragSrc && dragSrc !== row) {
      const list = document.getElementById('fields-list');
      const rows = [...list.children];
      const srcIdx = rows.indexOf(dragSrc);
      const tgtIdx = rows.indexOf(row);
      if (srcIdx < tgtIdx) {
        list.insertBefore(dragSrc, row.nextSibling);
      } else {
        list.insertBefore(dragSrc, row);
      }
      updatePreview();
    }
  });
}

function handleImageUpload(input) {
  const file = input.files[0];
  if (!file) return;
  uploadedImage = file;
  importedAscii = null;
  document.getElementById('upload-status').textContent = `✓ ${file.name}`;
  document.getElementById('upload-label').textContent = file.name;
  updatePreview();
}

function decodeConfig(b64) {
  const bytes = atob(b64);
  const json = decodeURIComponent(
    bytes.split('').map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)).join('')
  );
  return JSON.parse(json);
}

function parseSvgConfig(text) {
  const match = text.match(/<metadata id="mefetch-config">([^<]*)<\/metadata>/);
  if (!match) return null;
  try {
    return decodeConfig(match[1]);
  } catch (e) {
    return null;
  }
}

function applyConfig(cfg) {
  document.getElementById('username').value = cfg.username || '';
  document.getElementById('hostname').value = cfg.hostname || '';
  if (cfg.background) document.getElementById('background').value = cfg.background;
  if (cfg.keycolor)   document.getElementById('keycolor').value   = cfg.keycolor;
  if (cfg.textcolor)  document.getElementById('textcolor').value  = cfg.textcolor;
  document.getElementById('showstats').checked = cfg.showstats !== false;

  const list = document.getElementById('fields-list');
  list.innerHTML = '';
  (cfg.fields || []).forEach(f => {
    if (f === '---space---') {
      addSpacer();
    } else if (f.startsWith('---:')) {
      addSection(f.slice(4));
    } else {
      const idx = f.indexOf(':');
      if (idx >= 0) addField(f.slice(0, idx), f.slice(idx + 1));
    }
  });

  uploadedImage = null;
  importedAscii = cfg.ascii || null;
  document.getElementById('image-upload').value = '';
  if (importedAscii) {
    document.getElementById('upload-status').textContent = '✓ ascii art restored from svg';
    document.getElementById('upload-label').textContent = 'ascii art restored';
  } else {
    document.getElementById('upload-status').textContent = '';
    document.getElementById('upload-label').textContent = 'click to upload image';
  }
  updatePreview();
}

function handleSvgImport(input) {
  const file = input.files[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = () => {
    const cfg = parseSvgConfig(reader.result);
    const status = document.getElementById('import-status');
    if (!cfg) {
      status.style.color = '#f85149';
      status.textContent = '✕ not a mefetch svg';
      return;
    }
    applyConfig(cfg);
    status.style.color = '#3fb950';
    status.textContent = `✓ ${file.name}`;
  };
  reader.readAsText(file);
}

function removeSvg() {
  uploadedImage = null;
  importedAscii = null;
  document.getElementById('image-upload').value = '';
  document.getElementById('svg-import').value = '';
  document.getElementById('upload-status').textContent = '';
  document.getElementById('upload-label').textContent = 'click to upload image';
  document.getElementById('import-status').textContent = '';
  updatePreview();
}

async function downloadCard() {
  const resp = await fetch('/card/generate', { method: 'POST', body: buildForm() });
  const blob = await resp.blob();
  const url  = URL.createObjectURL(blob);
  const a    = document.createElement('a');
  a.href     = url;
  a.download = 'mefetch-card.svg';
  a.click();
  URL.revokeObjectURL(url);
}