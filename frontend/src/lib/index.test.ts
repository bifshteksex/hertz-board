import { describe, it, expect } from 'vitest';
import { APP_NAME } from './index';

describe('lib exports', () => {
	it('should export APP_NAME', () => {
		expect(APP_NAME).toBe('HertzBoard');
	});
});
