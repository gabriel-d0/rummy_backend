import { describe, it, expect } from 'vitest';
import {
	Version,
	OpClientStart,
	OpClientDiscard,
	OpServerState,
	OpServerError,
	NewEnvelope,
	parseEnvelope
} from './protocol';

describe('Nakama envelope — Day 18', () => {
	it('Version is 1', () => {
		expect(Version).toBe(1);
	});

	it('opcodes are stable 1..9 and 100..103', () => {
		expect(OpClientStart).toBe(1);
		expect(OpClientDiscard).toBe(2);
		expect(OpServerState).toBe(100);
		expect(OpServerError).toBe(102);
	});

	it('NewEnvelope creates Version 1 envelope with op and payload', () => {
		const json = NewEnvelope(1, { foo: 'bar' }, 'req-1');
		const obj = JSON.parse(json);
		expect(obj.v).toBe(1);
		expect(obj.op).toBe(1);
		expect(obj.requestId).toBe('req-1');
		expect(obj.payload).toEqual({ foo: 'bar' });
	});

	it('NewEnvelope without requestId', () => {
		const json = NewEnvelope(OpClientDiscard, { tileId: 't1' });
		const obj = JSON.parse(json);
		expect(obj.v).toBe(1);
		expect(obj.op).toBe(OpClientDiscard);
		expect(obj.payload).toEqual({ tileId: 't1' });
	});

	it('parseEnvelope parses and validates', () => {
		const json = NewEnvelope(OpClientDiscard, { tileId: 't1' }, 'req-2');
		const env = parseEnvelope(json);
		expect(env.v).toBe(1);
		expect(env.op).toBe(OpClientDiscard);
		expect(env.requestId).toBe('req-2');
	});

	it('parseEnvelope throws on bad envelope', () => {
		expect(() => parseEnvelope(JSON.stringify({ op: 1 }))).toThrow();
		expect(() => parseEnvelope(JSON.stringify({ v: 1 }))).toThrow();
	});
});
