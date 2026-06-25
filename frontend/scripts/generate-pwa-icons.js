// 生成 PWA 图标脚本
// 使用 Node.js 内置的 fs 和简单的画布绘制来生成图标

const fs = require('fs');
const path = require('path');

const sourceIcon = path.join(__dirname, '../../build/appicon.png');
const publicDir = path.join(__dirname, '../public');

// 确保 public 目录存在
if (!fs.existsSync(publicDir)) {
    fs.mkdirSync(publicDir, { recursive: true });
}

console.log('正在复制源图标作为基础...');

// 由于没有 ImageMagick 或 sharp，我们先将源图标复制过去
// 在实际部署时，用户需要使用在线工具或其他设备生成正确尺寸的图标

const iconSizes = [
    { name: 'pwa-192x192.png', size: 192 },
    { name: 'pwa-512x512.png', size: 512 },
    { name: 'apple-touch-icon-180x180.png', size: 180 },
    { name: 'favicon.ico', size: 32 }
];

// 先复制源图标到 public 目录
fs.copyFileSync(sourceIcon, path.join(publicDir, 'appicon-source.png'));

console.log(`
⚠️  图标生成需要 ImageMagick 或 sharp 库

由于本地环境限制，请使用以下方式之一生成 PWA 图标：

方式 1 - 在线工具（推荐）：
访问 https://realfavicongenerator.net/
上传 build/appicon.png，下载生成的图标包，解压到 frontend/public/

方式 2 - 使用 ImageMagick（如果其他设备有）：
cd frontend/public
convert ../../build/appicon.png -resize 192x192 pwa-192x192.png
convert ../../build/appicon.png -resize 512x512 pwa-512x512.png
convert ../../build/appicon.png -resize 180x180 apple-touch-icon-180x180.png
convert ../../build/appicon.png -resize 32x32 favicon.ico

方式 3 - 临时方案（开发测试用）：
将源图标复制为各个尺寸（浏览器会自动缩放，但不完美）：
`);

// 临时方案：复制源图标作为各个尺寸（仅用于开发测试）
iconSizes.forEach(({ name }) => {
    const targetPath = path.join(publicDir, name);
    fs.copyFileSync(sourceIcon, targetPath);
    console.log(`✓ 已创建 ${name} (临时使用源图标)`);
});

console.log(`
✓ 临时图标已创建在 frontend/public/
  这些图标在开发测试时可用，但生产环境建议使用正确尺寸的图标。
`);
