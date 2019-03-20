const explorerApi = require('../api/v1/ExplorerApi');
const openapi = require('../openapi');
const Ajv = require('ajv'), ajv = new Ajv();
ajv.addMetaSchema(require('ajv/lib/refs/json-schema-draft-06.json'));
ajv.addSchema(openapi, 'openapi.json');

describe('statistics api methods', () => {
    it('get statistics', () => {
        return explorerApi.getStatistics().then(response => {
            expect(response.status).toEqual(200);
            expect(ajv.validate({$ref: 'openapi.json#/components/schemas/StatisticCollectionResponse'}, response.data)).toEqual(true);
        });
    });

    it('get statistics by time filter', () => {
        return explorerApi.getStatistics().then(baseResponse => {
            startTime = baseResponse.data.data[0].date;
            endTime = baseResponse.data.data[1].date;
            explorerApi.getStatistics(null, startTime, endTime).then(response => {
                expect(response.status).toEqual(200);
                expect(response.data.data.length).toEqual(1);
                expect(ajv.validate({$ref: 'openapi.json#/components/schemas/StatisticCollectionResponse'}, response.data)).toEqual(true);
            });
        });
    });
});