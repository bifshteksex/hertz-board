# Contributing to HertzBoard

[English](#english) | [Русский](#russian) | [中文](#chinese)

---

<a name="english"></a>

## 🇬🇧 English

Thank you for considering contributing to HertzBoard! We welcome contributions from everyone.

### Code of Conduct

Please be respectful and constructive in all interactions.

### How to Contribute

#### Reporting Bugs

1. Check if the bug has already been reported in [Issues](https://github.com/bifshteksex/hertz-board/issues)
2. If not, create a new issue using the bug report template
3. Include:
   - Clear description of the bug
   - Steps to reproduce
   - Expected vs actual behavior
   - Screenshots if applicable
   - Environment details (OS, browser, version)

#### Suggesting Features

1. Check if the feature has already been suggested
2. Create a new issue using the feature request template
3. Describe:
   - The problem you're trying to solve
   - Your proposed solution
   - Alternative solutions you've considered
   - Use cases

#### Pull Requests

1. Fork the repository
2. Create a new branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Make your changes
4. Write or update tests
5. Ensure all tests pass
6. Follow the code style guidelines
7. Commit your changes following [Conventional Commits](https://www.conventionalcommits.org/):
   ```bash
   git commit -m "feat: add new feature"
   ```
8. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```
9. Create a Pull Request to the `main` branch

### Development Setup

See [Development Setup Guide](docs/development/setup.md) for detailed instructions.

Quick start:
```bash
git clone https://github.com/bifshteksex/hertz-board.git
cd hertz-board
make init
```

#### Commit Convention Setup

This project enforces [Conventional Commits](https://www.conventionalcommits.org/) at multiple levels:

**Local validation (Husky + Commitlint):**
- Automatically installed when you run `npm install` in the root directory
- Validates commit messages before they are created
- Prevents non-compliant commits on your local machine

**GitHub Actions:**
- Validates all commits in Pull Requests
- Checks both individual commits and PR title
- PRs with invalid commit messages will fail CI checks

**To test your commit message format:**
```bash
# This will validate your message before committing
git commit -m "feat(canvas): add new drawing tool"
```

If your commit message doesn't follow the convention, you'll see an error like:
```
⧗   input: invalid commit message
✖   subject may not be empty [subject-empty]
✖   type may not be empty [type-empty]
```

### Code Style

#### Backend (Go)

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Run `golangci-lint` before committing
- Write tests for new code
- Keep functions small and focused

#### Frontend (TypeScript/Svelte)

- Use TypeScript for type safety
- Follow the Prettier configuration
- Use meaningful variable names
- Write unit tests for components
- Keep components focused on a single responsibility

### Commit Messages

We use [Conventional Commits](https://www.conventionalcommits.org/) format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types (strictly enforced):**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks
- `perf`: Performance improvements
- `ci`: CI/CD changes
- `build`: Build system changes
- `revert`: Revert previous commit

**Rules:**
- Type must be lowercase
- Subject cannot be empty
- Subject cannot end with a period
- Header max length: 100 characters
- Body and footer must have blank line before them

**Valid examples:**
```
feat(canvas): add shape rotation feature
fix(auth): resolve JWT token expiration issue
docs(api): update REST API documentation
perf(renderer): optimize canvas rendering performance
ci(actions): add commit message validation
```

**Invalid examples:**
```
Feature: add rotation        ❌ Type must be from allowed list
feat(Canvas): Add rotation   ❌ Type and scope must be lowercase
feat: add rotation.          ❌ Subject cannot end with period
feat:add rotation            ❌ Missing space after colon
add rotation                 ❌ Missing type
```

### Testing

#### Backend Tests

```bash
make backend-test
```

#### Frontend Tests

```bash
make frontend-test
```

#### All Tests

```bash
make test
```

### Pull Request Checklist

- [ ] Code follows the project's style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] No merge conflicts

### Review Process

1. At least one maintainer will review your PR
2. Address any feedback or requested changes
3. Once approved, a maintainer will merge your PR

### Questions?

Feel free to ask questions by:
- Opening an issue
- Commenting on an existing issue or PR
- Reaching out to the maintainers

### License

By contributing, you agree that your contributions will be licensed under the GPL-3.0 License.

---

<a name="russian"></a>

## 🇷🇺 Русский

Спасибо за ваш интерес к HertzBoard! Мы приветствуем вклад каждого.

### Кодекс поведения

Пожалуйста, будьте уважительны и конструктивны во всех взаимодействиях.

### Как внести вклад

#### Сообщение об ошибках

1. Проверьте, не была ли ошибка уже сообщена в [Issues](https://github.com/bifshteksex/hertz-board/issues)
2. Если нет, создайте новую issue по шаблону bug report
3. Укажите:
   - Четкое описание ошибки
   - Шаги воспроизведения
   - Ожидаемое и фактическое поведение
   - Скриншоты (если применимо)
   - Детали окружения (ОС, браузер, версия)

#### Предложение новых функций

1. Проверьте, не была ли функция уже предложена
2. Создайте новую issue по шаблону feature request
3. Опишите:
   - Проблему, которую вы пытаетесь решить
   - Ваше предложенное решение
   - Альтернативные решения, которые вы рассматривали
   - Случаи использования

#### Pull Request'ы

1. Сделайте fork репозитория
2. Создайте новую ветку от `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. Внесите изменения
4. Напишите или обновите тесты
5. Убедитесь, что все тесты проходят
6. Следуйте руководству по стилю кода
7. Зафиксируйте изменения, следуя [Conventional Commits](https://www.conventionalcommits.org/):
   ```bash
   git commit -m "feat: add new feature"
   ```
8. Отправьте в ваш fork:
   ```bash
   git push origin feature/your-feature-name
   ```
9. Создайте Pull Request в ветку `main`

### Настройка окружения разработки

См. [Руководство по настройке](docs/development/setup.md) для подробных инструкций.

Быстрый старт:
```bash
git clone https://github.com/bifshteksex/hertz-board.git
cd hertz-board
make init
```

#### Настройка соглашений о коммитах

Проект **строго требует** соблюдения [Conventional Commits](https://www.conventionalcommits.org/) на нескольких уровнях:

**Локальная валидация (Husky + Commitlint):**
- Автоматически устанавливается при запуске `npm install` в корневой директории
- Проверяет сообщения коммитов перед их созданием
- Предотвращает несоответствующие коммиты на вашей локальной машине

**GitHub Actions:**
- Проверяет все коммиты в Pull Request'ах
- Проверяет как отдельные коммиты, так и заголовок PR
- PR с неверными сообщениями коммитов не пройдут проверки CI

**Тестирование формата сообщения коммита:**
```bash
# Это проверит ваше сообщение перед коммитом
git commit -m "feat(canvas): add new drawing tool"
```

Если сообщение коммита не соответствует соглашению, вы увидите ошибку:
```
⧗   input: invalid commit message
✖   subject may not be empty [subject-empty]
✖   type may not be empty [type-empty]
```

### Стиль кода

#### Backend (Go)

- Следуйте [Effective Go](https://golang.org/doc/effective_go.html)
- Используйте `gofmt` для форматирования
- Запускайте `golangci-lint` перед коммитом
- Пишите тесты для нового кода
- Держите функции маленькими и сфокусированными

#### Frontend (TypeScript/Svelte)

- Используйте TypeScript для типобезопасности
- Следуйте конфигурации Prettier
- Используйте осмысленные имена переменных
- Пишите unit-тесты для компонентов
- Держите компоненты сфокусированными на одной ответственности

### Сообщения коммитов

Мы используем формат [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Типы (строго соблюдаются):**
- `feat`: Новая функциональность
- `fix`: Исправление ошибки
- `docs`: Изменения в документации
- `style`: Изменения стиля кода (форматирование и т.д.)
- `refactor`: Рефакторинг кода
- `test`: Добавление или обновление тестов
- `chore`: Рутинные задачи
- `perf`: Улучшение производительности
- `ci`: Изменения CI/CD
- `build`: Изменения системы сборки
- `revert`: Отмена предыдущего коммита

**Правила:**
- Тип должен быть в нижнем регистре
- Тема не может быть пустой
- Тема не может заканчиваться точкой
- Максимальная длина заголовка: 100 символов
- Тело и футер должны иметь пустую строку перед собой

**Правильные примеры:**
```
feat(canvas): add shape rotation feature
fix(auth): resolve JWT token expiration issue
docs(api): update REST API documentation
perf(renderer): optimize canvas rendering performance
ci(actions): add commit message validation
```

**Неправильные примеры:**
```
Feature: add rotation        ❌ Тип должен быть из разрешенного списка
feat(Canvas): Add rotation   ❌ Тип и область должны быть в нижнем регистре
feat: add rotation.          ❌ Тема не может заканчиваться точкой
feat:add rotation            ❌ Отсутствует пробел после двоеточия
add rotation                 ❌ Отсутствует тип
```

### Тестирование

#### Тесты Backend

```bash
make backend-test
```

#### Тесты Frontend

```bash
make frontend-test
```

#### Все тесты

```bash
make test
```

### Чеклист Pull Request

- [ ] Код соответствует руководству по стилю проекта
- [ ] Выполнен самостоятельный code review
- [ ] Добавлены комментарии для сложного кода
- [ ] Обновлена документация
- [ ] Нет новых предупреждений
- [ ] Тесты добавлены/обновлены
- [ ] Все тесты проходят
- [ ] Нет конфликтов слияния

### Процесс ревью

1. Как минимум один мейнтейнер проверит ваш PR
2. Устраните все замечания и запрошенные изменения
3. После одобрения мейнтейнер сольёт ваш PR

### Вопросы?

Не стесняйтесь задавать вопросы:
- Создав issue
- Комментируя существующую issue или PR
- Связавшись с мейнтейнерами

### Лицензия

Внося вклад, вы соглашаетесь, что ваши изменения будут лицензированы под лицензией GPL-3.0.

---

<a name="chinese"></a>

## 🇨🇳 中文

感谢您考虑为 HertzBoard 做出贡献！我们欢迎所有人的贡献。

### 行为准则

请在所有互动中保持尊重和建设性。

### 如何贡献

#### 报告错误

1. 在 [Issues](https://github.com/bifshteksex/hertz-board/issues) 中检查错误是否已被报告
2. 如果没有，使用错误报告模板创建新问题
3. 包括：
   - 错误的清晰描述
   - 重现步骤
   - 期望行为与实际行为
   - 屏幕截图（如适用）
   - 环境详情（操作系统、浏览器、版本）

#### 建议新功能

1. 检查该功能是否已被建议
2. 使用功能请求模板创建新问题
3. 描述：
   - 您要解决的问题
   - 您提议的解决方案
   - 您考虑过的替代方案
   - 用例

#### Pull Request

1. Fork 仓库
2. 从 `main` 创建新分支：
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. 进行更改
4. 编写或更新测试
5. 确保所有测试通过
6. 遵循代码风格指南
7. 按照 [Conventional Commits](https://www.conventionalcommits.org/) 提交更改：
   ```bash
   git commit -m "feat: add new feature"
   ```
8. 推送到您的 fork：
   ```bash
   git push origin feature/your-feature-name
   ```
9. 创建 Pull Request 到 `main` 分支

### 开发环境设置

详细说明请参见[开发设置指南](docs/development/setup.md)。

快速开始：
```bash
git clone https://github.com/bifshteksex/hertz-board.git
cd hertz-board
make init
```

#### 提交约定设置

本项目在多个级别**严格执行** [Conventional Commits](https://www.conventionalcommits.org/)：

**本地验证（Husky + Commitlint）：**
- 在根目录运行 `npm install` 时自动安装
- 在创建提交消息之前验证
- 防止在本地机器上进行不合规的提交

**GitHub Actions：**
- 验证 Pull Request 中的所有提交
- 检查单个提交和 PR 标题
- 具有无效提交消息的 PR 将无法通过 CI 检查

**测试提交消息格式：**
```bash
# 这将在提交前验证您的消息
git commit -m "feat(canvas): add new drawing tool"
```

如果提交消息不符合约定，您将看到如下错误：
```
⧗   input: invalid commit message
✖   subject may not be empty [subject-empty]
✖   type may not be empty [type-empty]
```

### 代码风格

#### 后端（Go）

- 遵循 [Effective Go](https://golang.org/doc/effective_go.html)
- 使用 `gofmt` 进行格式化
- 提交前运行 `golangci-lint`
- 为新代码编写测试
- 保持函数小而专注

#### 前端（TypeScript/Svelte）

- 使用 TypeScript 实现类型安全
- 遵循 Prettier 配置
- 使用有意义的变量名
- 为组件编写单元测试
- 保持组件专注于单一职责

### 提交消息

我们使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<type>(<scope>): <subject>

<body>

<footer>
```

**类型（严格执行）：**
- `feat`: 新功能
- `fix`: 错误修复
- `docs`: 文档更改
- `style`: 代码样式更改（格式化等）
- `refactor`: 代码重构
- `test`: 添加或更新测试
- `chore`: 维护任务
- `perf`: 性能改进
- `ci`: CI/CD 更改
- `build`: 构建系统更改
- `revert`: 还原之前的提交

**规则：**
- 类型必须小写
- 主题不能为空
- 主题不能以句点结尾
- 标题最大长度：100 个字符
- 正文和页脚之前必须有空行

**有效示例：**
```
feat(canvas): add shape rotation feature
fix(auth): resolve JWT token expiration issue
docs(api): update REST API documentation
perf(renderer): optimize canvas rendering performance
ci(actions): add commit message validation
```

**无效示例：**
```
Feature: add rotation        ❌ 类型必须来自允许列表
feat(Canvas): Add rotation   ❌ 类型和范围必须小写
feat: add rotation.          ❌ 主题不能以句点结尾
feat:add rotation            ❌ 冒号后缺少空格
add rotation                 ❌ 缺少类型
```

### 测试

#### 后端测试

```bash
make backend-test
```

#### 前端测试

```bash
make frontend-test
```

#### 所有测试

```bash
make test
```

### Pull Request 检查清单

- [ ] 代码遵循项目的风格指南
- [ ] 完成自我审查
- [ ] 为复杂代码添加注释
- [ ] 更新文档
- [ ] 没有生成新的警告
- [ ] 添加/更新测试
- [ ] 所有测试通过
- [ ] 没有合并冲突

### 审查流程

1. 至少一名维护者将审查您的 PR
2. 处理任何反馈或请求的更改
3. 获得批准后，维护者将合并您的 PR

### 有疑问？

随时提问：
- 创建 issue
- 在现有 issue 或 PR 上评论
- 联系维护者

### 许可证

通过贡献，您同意您的贡献将根据 GPL-3.0 许可证获得许可。
