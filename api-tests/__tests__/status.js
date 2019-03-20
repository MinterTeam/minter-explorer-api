const explorerApi = require('../api/v1/ExplorerApi');
const openapi = require('../openapi');
const Ajv = require('ajv'), ajv = new Ajv();
ajv.addMetaSchema(require('ajv/lib/refs/json-schema-draft-06.json'));
ajv.addSchema(openapi, 'openapi.json');

describe('status api methods', () => {
    it('get status', () => {
        return explorerApi.getStatus().then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/StatusResponse'}, response.data)).toEqual(true);
        });
    });

    it('get status page', () => {
        return explorerApi.getStatusPage().then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/StatusPageResponse'}, response.data)).toEqual(true);
        });
    });
});