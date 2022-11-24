import * as wasmlib from '../lib/index';

describe('wasmlib', function () {

    describe('test', function () {
        it('should create service', () => {
            const client = new wasmlib.ScSandbox();
            // expect(client.Err() == null).toBeTruthy();
        });
    });
});
