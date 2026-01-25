# Branch Protection Rules

Для обеспечения качества кода и соблюдения Conventional Commits рекомендуется настроить Branch Protection Rules в GitHub.

## Настройка для `main` и `develop` веток

### Шаг 1: Перейдите в Settings

1. Откройте репозиторий на GitHub
2. Перейдите в **Settings** → **Branches**
3. Нажмите **Add branch protection rule**

### Шаг 2: Настройте правила для `main`

**Branch name pattern:** `main`

**Рекомендуемые настройки:**

✅ **Require a pull request before merging**
   - Require approvals: 1
   - Dismiss stale pull request approvals when new commits are pushed
   - Require review from Code Owners (опционально)

✅ **Require status checks to pass before merging**
   - Require branches to be up to date before merging
   - Status checks that are required:
     - `Validate Commit Messages` (из commitlint.yml)
     - `Backend Lint` (из ci.yml)
     - `Frontend Lint` (из ci.yml)
     - `Backend Tests` (из ci.yml)
     - `Frontend Tests` (из ci.yml)

✅ **Require conversation resolution before merging**
   - Требует разрешения всех комментариев перед мержем

✅ **Require linear history**
   - Запрещает merge commits, требует rebase или squash

✅ **Do not allow bypassing the above settings**
   - Даже администраторы должны следовать правилам

❌ **Allow force pushes** - ВЫКЛЮЧЕНО
❌ **Allow deletions** - ВЫКЛЮЧЕНО

### Шаг 3: Настройте правила для `develop`

Создайте аналогичное правило для ветки `develop` с теми же настройками, но можно ослабить требования:

**Branch name pattern:** `develop`

- Require approvals: 1 (можно оставить 0 для быстрой разработки)
- Все остальные настройки аналогично `main`

### Шаг 4: Настройте правила для feature веток (опционально)

**Branch name pattern:** `feature/*`

Более мягкие правила для feature веток:
- Не требовать approvals
- Требовать только прохождение CI checks
- Разрешить force push (для rebase)

## Результат

После настройки:

1. **Прямые коммиты в `main` и `develop` запрещены**
   - Все изменения только через Pull Requests

2. **Неправильные commit messages блокируются**
   - GitHub Action `Validate Commit Messages` должен пройти успешно
   - PR с неправильными коммитами не может быть смержен

3. **Все тесты и линтеры должны пройти**
   - Backend и Frontend линтеры
   - Все тесты
   - Build процесс

4. **Требуется code review**
   - Минимум 1 approval от другого разработчика

## Примеры блокировок

### ❌ Блокировано: Неправильный commit message

```
Commit: "add new feature"
Status: ❌ Validate Commit Messages failed
Reason: Type must be specified (feat, fix, docs, etc.)
```

### ❌ Блокировано: Тесты не проходят

```
Commit: "feat(auth): add login endpoint"
Status: ❌ Backend Tests failed
Reason: Unit tests are failing
```

### ✅ Разрешено: Все проверки прошли

```
Commit: "feat(auth): add login endpoint"
Status: ✅ All checks passed
- ✅ Validate Commit Messages
- ✅ Backend Lint
- ✅ Frontend Lint
- ✅ Backend Tests
- ✅ Frontend Tests
```

## Обход правил в экстренных случаях

Если необходимо срочно внести изменения:

1. **Временно отключите Branch Protection**
   - Settings → Branches → Edit rule
   - Снимите галочку с "Do not allow bypassing"

2. **Внесите изменения**

3. **Верните настройки обратно**

⚠️ **Внимание:** Использовать только в критических ситуациях!

## Проверка настроек

Чтобы убедиться, что правила работают:

1. Попробуйте сделать commit с неправильным сообщением:
```bash
git commit -m "wrong message"
```
Должна появиться ошибка от commitlint

2. Создайте PR с неправильным коммитом
   - GitHub Action должен зафейлиться
   - Merge button будет недоступен

3. Попробуйте сделать прямой push в `main`:
```bash
git push origin main
```
Должна появиться ошибка: "protected branch update failed"

## Дополнительные ресурсы

- [GitHub Branch Protection Documentation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Commitlint Documentation](https://commitlint.js.org/)
