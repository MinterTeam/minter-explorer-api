require('dotenv').config();
const explorerApi = require('../api/v1/ExplorerApi');
const openapi = require('../openapi');
const Ajv = require('ajv'), ajv = new Ajv();
ajv.addMetaSchema(require('ajv/lib/refs/json-schema-draft-06.json'));
ajv.addSchema(openapi, 'openapi.json');

const BASE_MINTER_ADDRESS = process.env.BASE_MINTER_ADDRESS;

describe('addresses api methods', () => {
    it('get addresses', () => {
        return explorerApi.getAddresses([BASE_MINTER_ADDRESS]).then(response => {
            expect(response.status).toEqual(200);
            expect(response.data.data.length).toEqual(1);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/AddressCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get addresses with empty params', () => {
        return explorerApi.getAddresses().catch(error => {
            expect(error.response.status).toEqual(422);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ValidationErrorResponse'}, error.response.data)).toEqual(true);
        });
    });

    it('get addresses with invalid params', () => {
        return explorerApi.getAddresses(['test']).catch(error => {
            expect(error.response.status).toEqual(422);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ValidationErrorResponse'}, error.response.data)).toEqual(true);
        });
    });

    it('get address', () => {
        return explorerApi.getAddress(BASE_MINTER_ADDRESS).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/AddressResponse'}, response.data.data)).toEqual(true);
        });
    });

    it('get address with invalid params', () => {
        return explorerApi.getAddress('test').catch(error => {
            expect(error.response.status).toEqual(422);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/ValidationErrorResponse'}, error.response.data)).toEqual(true);
        });
    });

    it('get address transactions', () => {
        return explorerApi.getAddressTransactions(BASE_MINTER_ADDRESS).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get address transactions with filter', () => {
        return explorerApi.getAddressTransactions(BASE_MINTER_ADDRESS).then(baseResponse => {
            blockId = baseResponse.data.data[0].block;
            explorerApi.getAddressTransactions(BASE_MINTER_ADDRESS, blockId, blockId).then(response => {
                expect(response.status).toEqual(200);
                expect(response.data.data.length).toEqual(1);
                expect(ajv.validate({$ref: 'openapi.json#/components/schemas/TransactionPaginatedCollectionResponse'}, response.data)).toEqual(true);
            });
        });
    });

    it('get address rewards', () => {
        return explorerApi.getAddressRewards(BASE_MINTER_ADDRESS).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/RewardPaginatedCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get address rewards with filter', () => {
        return explorerApi.getAddressRewards(BASE_MINTER_ADDRESS).then(baseResponse => {
            blockId = baseResponse.data.data[0].block;
            rowsCount = 0;
            
            baseResponse.data.data.forEach(function(item) {
                if (item.block == blockId) rowsCount++;
            });

            explorerApi.getAddressRewards(BASE_MINTER_ADDRESS, blockId, blockId).then(response => {
                expect(response.status).toEqual(200);
                expect(response.data.data.length).toEqual(rowsCount);
                expect(ajv.validate({$ref: 'openapi.json#/components/schemas/RewardPaginatedCollectionResponse'}, response.data)).toEqual(true);
            });
        });
    });

    it('get address slashes', () => {
        return explorerApi.getAddressSlashes(BASE_MINTER_ADDRESS).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/SlashPaginatedCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get address slashes with filter', () => {
        return explorerApi.getAddressSlashes(BASE_MINTER_ADDRESS).then(baseResponse => {
            blockId = baseResponse.data.data[0].block;
            explorerApi.getAddressSlashes(BASE_MINTER_ADDRESS, blockId, blockId).then(response => {
                expect(response.status).toEqual(200);
                expect(response.data.data.length).toEqual(1);
                expect(ajv.validate({$ref: 'openapi.json#/components/schemas/SlashPaginatedCollectionResponse'}, response.data)).toEqual(true);
            });
        });
    });

    it('get address rewards statistics', () => {
        return explorerApi.getAddressRewardsStats(BASE_MINTER_ADDRESS).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/RewardStatisticCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get address delegations', () => {
        return explorerApi.getAddressDelegations(BASE_MINTER_ADDRESS).then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/DelegationCollectionResponse'}, response.data)).toEqual(true);
        });
    });
});