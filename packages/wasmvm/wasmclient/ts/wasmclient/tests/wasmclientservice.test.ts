import {WasmClientService} from '../lib/wasmclientservice';

describe('wasmclient', function () {

    describe('Create service', function () {
        it('should create service', async () => {
            const client = WasmClientService.DefaultWasmClientService();
            expect(client.Err() == null).toBeTruthy();
        });
    });
});
