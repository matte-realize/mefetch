let uploadedImage = null;
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
  if (uploadedImage) form.append('image', uploadedImage);
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
    <input type="text" class="field-label" placeholder="label" value="${label}" oninput="updatePreview()">
    <input type="text" class="field-value" placeholder="value" value="${value}" oninput="updatePreview()">
    <button class="btn danger" onclick="removeField(this)">✕</button>
  `;
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
    <input type="text" class="field-label" placeholder="section name" value="${label}"
           oninput="updatePreview()" style="grid-column: span 2;">
    <input type="hidden" class="field-value" value="__divider__">
    <button class="btn danger" onclick="removeField(this)">✕</button>
  `;
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
  document.getElementById('upload-status').textContent = `✓ ${file.name}`;
  document.getElementById('upload-label').textContent = file.name;
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