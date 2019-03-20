require('dotenv').config();
const explorerApi = require('../api/v1/ExplorerApi');
const openapi = require('../openapi');
const Ajv = require('ajv'), ajv = new Ajv();
ajv.addMetaSchema(require('ajv/lib/refs/json-schema-draft-06.json'));
ajv.addSchema(openapi, 'openapi.json');

const BASE_VALIDATOR_PUB_KEY = process.env.BASE_VALIDATOR_PUB_KEY;

describe('validators api methods', () => {
    it('get validator', () => {
        return explorerApi.getValidator(BASE_VALIDATOR_PUB_KEY).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ValidatorResponse'}, response.data)).toEqual(true);
        });
    });

    it('get unknown validator', () => {
        unknownPubKey = 'Mpc8c6834da8ba2b0b24f7e5ab67049509278e709cde925f14184586f74dcc9d0a';
        return explorerApi.getValidator(unknownPubKey).catch(error => {
            expect(error.response.status).toEqual(404);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ErrorResponse'}, error.response.data)).toEqual(true);
        });
    });

    it('get validator with invalid pub key', () => {
        return explorerApi.getValidator('test').catch(error => {
            expect(error.response.status).toEqual(422);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ValidationErrorResponse'}, error.response.data)).toEqual(true);
        });
    });

    it('get validator transactions', () => {
        return explorerApi.getValidatorTransactions(BASE_VALIDATOR_PUB_KEY).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get validator transactions with range filter', () => {
        return explorerApi.getValidatorTransactions(BASE_VALIDATOR_PUB_KEY).then(baseResponse => {
            blockId = baseResponse.data.data[0].block;
            explorerApi.getValidatorTransactions(BASE_VALIDATOR_PUB_KEY, blockId, blockId).then(response => {
                expect(response.status).toEqual(200);
                expect(response.data.data.length).toEqual(1);
                expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
            }); 
        });
    });
});