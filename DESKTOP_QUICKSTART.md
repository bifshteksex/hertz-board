# HertzBoard Desktop - Quick Start Guide

## 🚀 Для разработчиков

### Запуск в режиме разработки

```bash
cd frontend
npm install
npm run tauri:dev
```

### Сборка для текущей платформы

```bash
npm run tauri:build
```

Результат: `frontend/src-tauri/target/release/bundle/`

---

## 📦 Создание релиза

### Автоматический релиз (Recommended)

1. Создайте git tag:
```bash
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0
```

2. GitHub Actions автоматически:
   - Соберёт приложение для Windows, macOS, Linux
   - Создаст GitHub Release
   - Загрузит все артефакты

### Результаты сборки

**Windows:**
- `HertzBoard_x64_en-US.msi` - MSI installer

**macOS:**
- `HertzBoard.dmg` - DMG disk image
- `HertzBoard_x64.app.tar.gz` - Intel bundle
- `HertzBoard_aarch64.app.tar.gz` - Apple Silicon bundle

**Linux:**
- `hertzboard_amd64.AppImage` - Portable AppImage
- `hertzboard_amd64.deb` - Debian package

---

## 🔧 Требования для локальной сборки

### Windows
- Node.js 20+
- Rust (latest stable)
- Visual Studio Build Tools

### macOS
- Node.js 20+
- Rust (latest stable)
- Xcode Command Line Tools

### Linux (Ubuntu/Debian)
```bash
sudo apt-get install -y \
  libwebkit2gtk-4.1-dev \
  libappindicator3-dev \
  librsvg2-dev \
  patchelf
```

---

## 📝 Настройка перед production

### 1. Обновите конфигурацию

Отредактируйте `frontend/src-tauri/tauri.conf.json`:

```json
{
  "identifier": "com.yourcompany.hertzboard",
  "bundle": {
    "publisher": "Your Company Name"
  },
  "plugins": {
    "updater": {
      "endpoints": [
        "https://github.com/YOUR_USERNAME/hertz-board/releases/latest/download/latest.json"
      ]
    }
  }
}
```

### 2. Создайте иконки

Замените иконки в `frontend/src-tauri/icons/`:
- `32x32.png`
- `128x128.png`
- `icon.icns` (macOS)
- `icon.ico` (Windows)

### 3. (Опционально) Настройте code signing

**Windows:**
```yaml
# .github/workflows/desktop-release.yml
env:
  WINDOWS_CERTIFICATE: ${{ secrets.WINDOWS_CERTIFICATE }}
```

**macOS:**
```bash
# Local build
export APPLE_CERTIFICATE="Developer ID Application: Your Name"
```

---

## 🧪 Тестирование

```bash
# TypeScript check
npm run check

# Frontend tests
npm run test:unit

# Lint
npm run lint

# Debug build (faster)
npm run tauri:build:debug
```

---

## 📖 Дополнительная информация

- [Полная документация](frontend/DESKTOP.md)
- [Детали фазы 12](PHASE_12_SUMMARY.md)
- [Tauri Docs](https://v2.tauri.app/)

---

## ❓ Troubleshooting

**Ошибка: "webkit2gtk not found"** (Linux)
```bash
sudo apt-get install libwebkit2gtk-4.1-dev
```

**Ошибка: "Rust version too old"**
```bash
rustup update
```

**Windows SmartScreen warning**
- Это нормально для unsigned приложений
- Решение: Настройте code signing

---

## 🎯 Следующие шаги

1. ✅ Создайте первый релиз через git tag
2. ⏳ Протестируйте на всех платформах
3. ⏳ Настройте code signing
4. ⏳ Создайте custom иконки
5. ⏳ Напишите user documentation

**Успехов!** 🚀
