require('dotenv').config();
const BaseApi = require('../BaseApi');

class ExplorerApi extends BaseApi {
    constructor(baseUrl = null) {
        baseUrl = (baseUrl == null) ? process.env.BASE_URL : baseUrl;
        baseUrl += '/v1/';
        super(baseUrl);
    }

    // BLOCKS GROUP
    getBlock(id) {
        return this.client.get(`blocks/${id}`);
    }

    getBlocks(page = 1, limit = null) {
        return this.client.get(`blocks`, {
            params: {
                page: page,
                limit: limit
            }
        });
    }

    getBlockTransactions(id) {
        return this.client.get(`blocks/${id}/transactions`);
    }

    // COINS GROUP
    getCoins(symbol = null) {
        return this.client.get(`coins`, {
            params: {
                symbol: symbol
            }
        });
    }

    // ADDRESSES GROUP
    getAddresses(addresses = null) {
        return this.client.get(`addresses`, {
            params: {
                addresses: addresses
            }
        });
    }

    getAddress(address) {
        return this.client.get(`addresses/${address}`);
    }

    getAddressTransactions(address, startblock = null, endblock = null) {
        return this.client.get(`addresses/${address}/transactions`, {params: prepareRangeFilter(startblock, endblock)});
    }

    getAddressRewards(address, startblock = null, endblock = null) {
        return this.client.get(`addresses/${address}/events/rewards`, {params: prepareRangeFilter(startblock, endblock)});
    }

    getAddressSlashes(address, startblock = null, endblock = null) {
        return this.client.get(`addresses/${address}/events/slashes`, {params: prepareRangeFilter(startblock, endblock)});
    }

    getAddressRewardsStats(address, scale = null, startTime = null, endTime = null) {
        return this.client.get(`addresses/${address}/statistics/rewards`, {
            params: {
                scale: scale,
                startTime: startTime,
                endTime: endTime
            }
        });
    }

    getAddressDelegations(address) {
        return this.client.get(`addresses/${address}/delegations`);
    }

    // TRANSACTIONS GROUP
    getTransactions(addresses = null, startblock = null, endblock = null) {
        params = prepareRangeFilter(startblock, endblock);
        params.addresses = addresses;

        return this.client.get(`transactions`, {params: params});
    }

    getTransaction(hash) {
        return this.client.get(`transactions/${hash}`);
    }

    // VALIDATORS GROUP
    getValidator(publicKey) {
        return this.client.get(`validators/${publicKey}`);
    }

    getValidatorTransactions(publicKey, startblock = null, endblock = null) {
        return this.client.get(`validators/${publicKey}/transactions`, {params: prepareRangeFilter(startblock, endblock)});
    }

    // STATISTICS GROUP
    getStatistics(scale = null, startTime = null, endTime = null) {
        return this.client.get(`statistics/transactions`, {
            params: {
                scale: scale,
                startTime: startTime,
                endTime: endTime
            }
        });
    }

    // STATUS GROUP
    getStatus() {
        return this.client.get(`status`);
    }

    getStatusPage() {
        return this.client.get(`status-page`);
    }
}

function prepareRangeFilter(startblock = null, endblock = null) {
    params = {}

    if (startblock) {
        params.startblock = startblock;
    }

    if (endblock) {
        params.endblock = endblock;
    }

    return params;
}

module.exports = new ExplorerApi();