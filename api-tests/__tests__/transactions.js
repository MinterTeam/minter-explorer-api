require('dotenv').config();
const explorerApi = require('../api/v1/ExplorerApi');
const openapi = require('../openapi');
const Ajv = require('ajv'), ajv = new Ajv();
ajv.addMetaSchema(require('ajv/lib/refs/json-schema-draft-06.json'));
ajv.addSchema(openapi, 'openapi.json');

const BASE_MINTER_ADDRESS = process.env.BASE_MINTER_ADDRESS;

describe('transactions api methods', () => {
    it('get transactions', () => {
        return explorerApi.getTransactions().then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get transactions by range filter', () => {
        return explorerApi.getTransactions().then(baseResponse => {
            rowsCount = 0;
            blockId = baseResponse.data.data[0].block;
            baseResponse.data.data.forEach(function(item) {
                if (item.block === blockId) rowsCount++;
            })

            explorerApi.getTransactions(null, blockId, blockId).then(response => {
                expect(response.status).toEqual(200);
                expect(response.data.data.length).toEqual(rowsCount);
                expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
            });
        });
    });

    it('get transactions', () => {
        return explorerApi.getTransactions([BASE_MINTER_ADDRESS]).then(baseResponse => {
            rowsCount = 0;
            blockId = baseResponse.data.data[0].block;
            baseResponse.data.data.forEach(function(item) {
                if (item.block === blockId) rowsCount++;
            })

            explorerApi.getTransactions([BASE_MINTER_ADDRESS], blockId, blockId).then(response => {
                expect(response.status).toEqual(200);
                expect(response.data.data.length).toEqual(rowsCount);
                expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
            });
        });
    });

    it('get transaction', () => {
        txHash = 'Mt5c8eca5e5cfc18aeb9012e0516e2bdb61d5bddcb819f42cdf1a6fa05df32f7f1';
        return explorerApi.getTransaction(txHash).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionResponse'}, response.data.data)).toEqual(true);
        });
    });

    it('get unknown transaction', () => {
        txHash = 'Mt5c8eca5e5cfc18aeb9012e0516e2bdb61d5bddcb819f42cdf1a6fa05df32f7f0';
        return explorerApi.getTransaction(txHash).catch(error => {
            expect(error.response.status).toEqual(404);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ErrorResponse'}, error.response.data)).toEqual(true);
        });
    });
});