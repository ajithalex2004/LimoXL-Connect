const axios = require('axios');

async function testFetch() {
    try {
        console.log("Attempting fetch to http://localhost:8080/api/partner/rfqs");
        const response = await axios.get('http://localhost:8080/api/partner/rfqs');
        console.log("Status:", response.status);
        console.log("Data:", JSON.stringify(response.data, null, 2));
    } catch (error) {
        console.error("Error fetching RFQs:", error.message);
        if (error.response) {
            console.error("Response Status:", error.response.status);
            console.error("Response Data:", error.response.data);
        }
    }
}

testFetch();
