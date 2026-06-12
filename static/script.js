let uploadedImage = null;

window.onload = () => {
  addField('OS', '');
  addField('Shell', '');
  addField('Editor', '');
  updatePreview();
};

function getFields() {
  const rows = document.querySelectorAll('.field-row');
  const fields = [];
  
  rows.forEach(row => {
    const label = row.querySelector('.field-label').value.trim();
    const value = row.querySelector('.field-value').value.trim();
    if (label && value) fields.push({ label, value });
  });
  
  return fields;
}

function buildURL() {
  const username   = document.getElementById('username').value.trim();
  const hostname   = document.getElementById('hostname').value.trim();
  const background = document.getElementById('background').value;
  const keycolor   = document.getElementById('keycolor').value;
  const textcolor  = document.getElementById('textcolor').value;
  const showStats  = document.getElementById('showstats').checked;
  const params     = new URLSearchParams();
  
  if (username)   params.set('username', username);
  if (hostname)   params.set('hostname', hostname);
  if (!showStats) params.set('showstats', 'false');

  params.set('background', background);
  params.set('keycolor', keycolor);
  params.set('textcolor', textcolor);

  getFields().forEach(f => {
    params.append('field', `${f.label}:${f.value}`);
  });

  return `/card.svg?${params.toString()}`;
}

function updatePreview() {
  const url = buildURL();
  
  document.getElementById('card-preview').src = url + '&t=' + Date.now();
  updateSnippet(url);
}

function updateSnippet(url) {
  const absolute = `http://localhost:8080${url}`;
  
  document.getElementById('snippet-output').value =
    `![Neofetch Card](${absolute})`;
}

function copySnippet(type) {
  const url = buildURL();
  const absolute = `http://localhost:8080${url}`;
  const text = type === 'markdown'
    ? `![Neofetch Card](${absolute})`
    : absolute;
  
    navigator.clipboard.writeText(text);
}

function addField(label = '', value = '') {
  const list = document.getElementById('fields-list');
  const row = document.createElement('div');
  
  row.className = 'field-row';
  row.innerHTML = `
    <input type="text" class="field-label" placeholder="label" value="${label}" oninput="updatePreview()">
    <input type="text" class="field-value" placeholder="value" oninput="updatePreview()">
    <button class="btn danger" onclick="removeField(this)">✕</button>
  `;
  
  list.appendChild(row);
  updatePreview();
}

function removeField(btn) {
  btn.closest('.field-row').remove();
  updatePreview();
}

function handleImageUpload(input) {
  const file = input.files[0];
  if (!file) return;
  uploadedImage = file;
  
  document.getElementById('upload-status').textContent = `✓ ${file.name}`;
  document.getElementById('upload-label').textContent = file.name;
}

async function downloadCard() {
  const form = new FormData();

  if (uploadedImage) {
    form.append('image', uploadedImage);
  }

  form.append('username', document.getElementById('username').value.trim());
  form.append('hostname', document.getElementById('hostname').value.trim());
  form.append('background', document.getElementById('background').value);
  form.append('keycolor', document.getElementById('keycolor').value);
  form.append('textcolor', document.getElementById('textcolor').value);

  getFields().forEach(f => {
    form.append('field', `${f.label}:${f.value}`);
  });

  const resp = await fetch('/card/generate', { method: 'POST', body: form });
  const blob = await resp.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  
  a.href = url;
  a.download = 'neofetch-card.svg';
  a.click();
  URL.revokeObjectURL(url);
}