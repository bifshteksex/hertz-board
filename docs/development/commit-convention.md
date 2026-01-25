# Commit Convention Guide

## Quick Start

HertzBoard строго требует соблюдения [Conventional Commits](https://www.conventionalcommits.org/) спецификации.

### Формат

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

### Примеры

✅ **Правильно:**
```bash
git commit -m "feat(canvas): add shape rotation feature"
git commit -m "fix(auth): resolve JWT token expiration issue"
git commit -m "docs(readme): update installation instructions"
git commit -m "perf(renderer): optimize canvas rendering"
git commit -m "ci(actions): add commit validation workflow"
```

❌ **Неправильно:**
```bash
git commit -m "add new feature"              # Нет типа
git commit -m "Feature: add rotation"        # Type не из списка
git commit -m "feat(Canvas): Add feature"    # Scope с большой буквы
git commit -m "feat: add feature."           # Точка в конце subject
git commit -m "feat:add feature"             # Нет пробела после :
```

## Допустимые типы

| Тип | Описание | Пример |
|-----|----------|--------|
| `feat` | Новая функциональность | `feat(auth): add OAuth login` |
| `fix` | Исправление бага | `fix(canvas): resolve rendering glitch` |
| `docs` | Изменения в документации | `docs(api): update REST endpoints` |
| `style` | Форматирование кода (не влияющее на логику) | `style(backend): format go files` |
| `refactor` | Рефакторинг кода | `refactor(store): simplify state management` |
| `perf` | Улучшение производительности | `perf(db): optimize query execution` |
| `test` | Добавление или обновление тестов | `test(auth): add login unit tests` |
| `build` | Изменения в системе сборки | `build(docker): update base image` |
| `ci` | Изменения в CI/CD | `ci(github): add linting workflow` |
| `chore` | Рутинные задачи | `chore(deps): update dependencies` |
| `revert` | Отмена предыдущего коммита | `revert: feat(auth): add OAuth login` |

## Правила

### Type
- **Обязателен**
- Должен быть из списка выше
- Только в нижнем регистре
- Без точки в конце

### Scope (опционально)
- Указывает на модуль/компонент проекта
- В нижнем регистре
- В скобках после типа

Примеры scope:
- `auth` - Аутентификация
- `canvas` - Canvas функциональность
- `api` - API endpoints
- `db` - База данных
- `ws` - WebSocket
- `ui` - UI компоненты
- `store` - State management

### Subject
- **Обязателен**
- Краткое описание изменений
- Начинается с маленькой буквы
- Без точки в конце
- Использовать повелительное наклонение ("add", а не "added" или "adds")

### Body (опционально)
- Детальное описание изменений
- Отделяется пустой строкой от subject
- Может содержать несколько параграфов

### Footer (опционально)
- Ссылки на issues: `Closes #123`
- Breaking changes: `BREAKING CHANGE: description`
- Отделяется пустой строкой от body

## Примеры с body и footer

```bash
git commit -m "feat(auth): add OAuth 2.0 authentication

Implemented OAuth 2.0 flow with support for:
- Google OAuth provider
- GitHub OAuth provider
- Token refresh mechanism

Closes #45"
```

```bash
git commit -m "fix(api)!: change response format

BREAKING CHANGE: API response format changed from array to object
with pagination metadata.

Migration guide:
- Update client code to expect object instead of array
- Access data via response.data instead of response

Closes #67"
```

## Автоматическая валидация

### Локально (Husky)

При попытке создать коммит с неправильным сообщением:

```bash
$ git commit -m "add new feature"

⧗   input: add new feature
✖   subject may not be empty [subject-empty]
✖   type may not be empty [type-empty]

✖   found 2 problems, 0 warnings
```

Коммит **не будет создан**.

### В Pull Requests (GitHub Actions)

При создании PR проверяются:
1. Все коммиты в PR
2. Название PR (тоже должно следовать Conventional Commits)

Если проверка не прошла, PR **не может быть смержен**.

## Настройка

Автоматическая валидация настраивается при первой установке зависимостей:

```bash
# В корне проекта
npm install
```

Это установит:
- `@commitlint/cli` - Валидатор коммитов
- `@commitlint/config-conventional` - Правила Conventional Commits
- `husky` - Git hooks manager

Git hook автоматически создается в `.husky/commit-msg`.

## Отключение валидации (не рекомендуется)

Если по какой-то причине нужно обойти валидацию:

```bash
git commit -m "your message" --no-verify
```

⚠️ **Внимание:** Не рекомендуется использовать `--no-verify`, так как это нарушает стандарты проекта и PR всё равно не пройдет GitHub Actions.

## Проверка существующих коммитов

Проверить последние N коммитов:

```bash
npx commitlint --from HEAD~5
```

Проверить коммиты в текущей ветке относительно main:

```bash
npx commitlint --from main
```

## FAQ

**Q: Можно ли использовать emoji в коммитах?**
A: Нет, это не соответствует Conventional Commits спецификации.

**Q: Что делать, если нужно сделать несколько изменений в одном коммите?**
A: Лучше разбить на несколько коммитов, каждый с одним типом изменений. Если невозможно, используйте самый значимый тип.

**Q: Как писать коммиты для исправлений в нескольких модулях?**
A: Либо разбить на несколько коммитов, либо опустить scope:
```bash
fix: resolve authentication issues across modules
```

**Q: Нужно ли следовать этому формату в личных feature ветках?**
A: Да, так как при создании PR все коммиты будут проверены.

**Q: Что если я забыл правильный формат?**
A: Просто сделайте коммит - если формат неправильный, вы увидите ошибку с подсказкой.

## Дополнительные ресурсы

- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [Commitlint Documentation](https://commitlint.js.org/)
- [Husky Documentation](https://typicode.github.io/husky/)
- [Angular Commit Guidelines](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#-commit-message-format) (на чем основан Conventional Commits)
