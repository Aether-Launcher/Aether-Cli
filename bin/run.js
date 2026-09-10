#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { spawnSync } = require('child_process');
const { getPlatformInfo, install } = require('./install');

async function main() {
  const { localBinaryName } = getPlatformInfo();
  const binaryPath = path.join(__dirname, localBinaryName);

  // If binary is missing, attempt to download it now
  if (!fs.existsSync(binaryPath)) {
    console.log('[Aether CLI] First-time setup: fetching native binary...');
    await install();
  }

  if (!fs.existsSync(binaryPath)) {
    console.error(`[Aether CLI] Error: Native executable '${localBinaryName}' could not be found.`);
    console.error('Please verify your internet connection or install Go to build from source:');
    console.error('  go install github.com/Aether-Launcher/aether-cli/...@latest');
    process.exit(1);
  }

  // Ensure executable permissions on Unix
  try {
    fs.chmodSync(binaryPath, 0o755);
  } catch {}

  const result = spawnSync(binaryPath, process.argv.slice(2), {
    stdio: 'inherit',
    windowsHide: false
  });

  if (result.error) {
    console.error('[Aether CLI] Execution error:', result.error.message);
    process.exit(1);
  }

  process.exit(result.status !== null ? result.status : 0);
}

main().catch((err) => {
  console.error('[Aether CLI] Fatal error:', err.message);
  process.exit(1);
});
