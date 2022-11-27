import * as wasmlib from '../index';

describe('wasmlib', function () {

    describe('test', function () {
        it('string conversion', () => {
            const testValue = "Some weird test string";
            const buf = wasmlib.stringToBytes(testValue);
            expect(wasmlib.stringFromBytes(buf) == testValue).toBeTruthy();
            const str = wasmlib.stringToString(testValue);
            expect(wasmlib.stringFromString(str) == testValue).toBeTruthy();
        });
        it('uint8 conversion', () => {
            const testValue = 123 as u8;
            const buf = wasmlib.uint8ToBytes(testValue);
            expect(wasmlib.uint8FromBytes(buf) == testValue).toBeTruthy();
            const str = wasmlib.uint8ToString(testValue);
            expect(wasmlib.uint8FromString(str) == testValue).toBeTruthy();
        });
        it('uint16 conversion', () => {
            const testValue = 12345 as u16;
            const buf = wasmlib.uint16ToBytes(testValue);
            expect(wasmlib.uint16FromBytes(buf) == testValue).toBeTruthy();
            const str = wasmlib.uint16ToString(testValue);
            expect(wasmlib.uint16FromString(str) == testValue).toBeTruthy();
        });
        it('uint32 conversion', () => {
            const testValue = 123456789 as u32;
            const buf = wasmlib.uint32ToBytes(testValue);
            expect(wasmlib.uint32FromBytes(buf) == testValue).toBeTruthy();
            const str = wasmlib.uint32ToString(testValue);
            expect(wasmlib.uint32FromString(str) == testValue).toBeTruthy();
        });
    });
});
