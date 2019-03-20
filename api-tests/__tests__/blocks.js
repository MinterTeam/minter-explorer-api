const explorerApi = require('../api/v1/ExplorerApi');
const openapi = require('../openapi');
const Ajv = require('ajv'), ajv = new Ajv();
ajv.addMetaSchema(require('ajv/lib/refs/json-schema-draft-06.json'));
ajv.addSchema(openapi, 'openapi.json');

describe('blocks api methods', () => {
    it('get block by id', () => {
        return explorerApi.getBlock(1).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/BlockResponse'}, response.data.data)).toEqual(true);
        });
    });

    it('get unknown block', () => {
        return explorerApi.getBlock(0).catch(error => {
            expect(error.response.status).toEqual(404);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ErrorResponse'}, error.response.data)).toEqual(true);
        });
    });

    it('get block with string as id', () => {
        return explorerApi.getBlock('test').catch(error => {
            expect(error.response.status).toEqual(422);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ValidationErrorResponse'}, error.response.data)).toEqual(true);
        });
    });

    it('get blocks', () => {
        return explorerApi.getBlocks().then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/BlockPaginatedCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get blocks with limit', () => {
        return explorerApi.getBlocks(1, 5).then(response => {
            expect(response.status).toEqual(200);
            expect(response.data.data.length).toEqual(5)
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/BlockPaginatedCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get block transactions', () => {
        return  explorerApi.getTransactions().then(txResponse => {
            if (txResponse.data.data.length > 0) {
                blockId = txResponse.data.data[0].block;
                explorerApi.getBlockTransactions(blockId).then(response => {
                    expect(response.status).toEqual(200);
                    expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
                });
            }
        })
    });
});