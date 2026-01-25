export default {
	extends: ['@commitlint/config-conventional'],
	rules: {
		'type-enum': [
			2,
			'always',
			[
				'feat', // New feature
				'fix', // Bug fix
				'docs', // Documentation changes
				'style', // Code style changes (formatting, etc.)
				'refactor', // Code refactoring
				'test', // Adding or updating tests
				'chore', // Maintenance tasks
				'perf', // Performance improvements
				'ci', // CI/CD changes
				'build', // Build system changes
				'revert' // Revert previous commit
			]
		],
		'type-case': [2, 'always', 'lower-case'],
		'type-empty': [2, 'never'],
		'scope-case': [2, 'always', 'lower-case'],
		'subject-empty': [2, 'never'],
		'subject-full-stop': [2, 'never', '.'],
		'header-max-length': [2, 'always', 100],
		'body-leading-blank': [2, 'always'],
		'footer-leading-blank': [2, 'always']
	}
};
