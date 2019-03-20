const axios = require('axios');

class BaseApi {
    constructor(baseUrl) {
        this.client = axios.create({
            baseURL: baseUrl,
            timeout: 6000,
            headers: {'Content-Type': 'application/json'}
        });
    }
}

module.exports = BaseApi;