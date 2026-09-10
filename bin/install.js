#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');
const https = require('https');

const packageJson = require('../package.json');
const VERSION = packageJson.version;
const REPO = 'Aether-Launcher/Aether-Cli';

function getPlatformInfo() {
  const platform = os.platform();
  const arch = os.arch();

  let osName = '';
  let ext = '';
  if (platform === 'win32') {
    osName = 'windows';
    ext = '.exe';
  } else if (platform === 'darwin') {
    osName = 'darwin';
  } else if (platform === 'linux') {
    osName = 'linux';
  } else {
    throw new Error(`Unsupported OS platform: ${platform}`);
  }

  let archName = '';
  if (arch === 'x64') {
    archName = 'amd64';
  } else if (arch === 'arm64') {
    archName = 'arm64';
  } else {
    throw new Error(`Unsupported architecture: ${arch}`);
  }

  const binaryFilename = `aether-${osName}-${archName}${ext}`;
  const localBinaryName = `aether${ext}`;
  return { osName, archName, binaryFilename, localBinaryName };
}

function downloadFile(url, dest, redirectCount = 0) {
  if (redirectCount > 5) {
    return Promise.reject(new Error('Too many redirects while downloading binary.'));
  }

  return new Promise((resolve, reject) => {
    https.get(url, { headers: { 'User-Agent': '@aethermc/cli' } }, (res) => {
      if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
        return resolve(downloadFile(res.headers.location, dest, redirectCount + 1));
      }

      if (res.statusCode !== 200) {
        return reject(new Error(`Failed to download binary: HTTP ${res.statusCode} (${url})`));
      }

      const file = fs.createWriteStream(dest);
      res.pipe(file);
      file.on('finish', () => {
        file.close(() => resolve());
      });
      file.on('error', (err) => {
        fs.unlink(dest, () => {});
        reject(err);
      });
    }).on('error', reject);
  });
}

async function install() {
  try {
    const { binaryFilename, localBinaryName } = getPlatformInfo();
    const binDir = path.join(__dirname);
    const targetPath = path.join(binDir, localBinaryName);

    // If local binary already exists (e.g. built locally), skip
    if (fs.existsSync(targetPath)) {
      try {
        fs.chmodSync(targetPath, 0o755);
      } catch {}
      return;
    }

    const downloadUrl = `https://github.com/${REPO}/releases/download/v${VERSION}/${binaryFilename}`;
    console.log(`[Aether CLI] Downloading native binary for ${os.platform()}-${os.arch()} (v${VERSION})...`);
    
    await downloadFile(downloadUrl, targetPath);
    try {
      fs.chmodSync(targetPath, 0o755);
    } catch {}

    console.log(`[Aether CLI] Successfully installed ${localBinaryName}.`);
  } catch (err) {
    // If download fails during npm install (e.g. offline, or version not published yet),
    // warn gracefully so npm install doesn't fail. run.js will retry on demand.
    console.warn(`[Aether CLI] Notice: Could not download native binary during install (${err.message}).`);
    console.warn('[Aether CLI] It will be downloaded automatically on first execution.');
  }
}

if (require.main === module) {
  install();
}

module.exports = { getPlatformInfo, downloadFile, install };
