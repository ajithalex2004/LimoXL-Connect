
const axios = require('axios');

async function runFlow() {
    const partnerClient = axios.create({ baseURL: 'http://localhost:8080/api/partner' });
    const operatorClient = axios.create({ baseURL: 'http://localhost:8080/api/operator' });

    try {
        // 0. Debug Check: List ALL Operator Trips
        console.log("0. DEBUG: Listing Operator Trips...");
        const opTrips = await operatorClient.get('/trips').then(r => r.data);
        const targetTrip = opTrips.find(t => t.booking_reference === 'RFQ-1004');
        console.log("   RFQ-1004 Status:", targetTrip ? targetTrip.status : "NOT FOUND");
        console.log("   RFQ-1004 ID:", targetTrip ? targetTrip.id : "N/A");

        let tripID = targetTrip ? targetTrip.id : null;

        if (targetTrip && targetTrip.status === 'MARKETPLACE_SEARCH') {
            // 1. List RFQs (Partner)
            console.log("1. Partner: Listing RFQs...");
            const response = await partnerClient.get('/rfqs');
            if (response.data && response.data.length > 0) {
                tripID = response.data[0].id;
                // 2. Submit Quote
                console.log(`2. Partner: Submitting Quote for ${tripID}...`);
                await partnerClient.post('/quotes', {
                    trip_id: tripID,
                    price: 150.00,
                    notes: "Best price for you"
                });
            }
        } else if (targetTrip && targetTrip.status === 'OFFERED') {
            console.log("   Trip is already OFFERED. Skipping Quote Submission.");
            tripID = targetTrip.id;
        }

        if (!tripID) {
            console.log("No usable trip found.");
            return;
        }

        // 3. List Quotes (Operator)
        console.log("3. Operator: Listing Quotes...");
        const quotes = await operatorClient.get('/quotes').then(r => r.data);
        const offer = quotes.find(q => q.trip_id === tripID && q.status === 'PENDING');

        if (!offer) {
            console.log("   No pending offer found for this trip. Checking if already accepted...");
            // If ALREADY accepted, it should be in step 5
        } else {
            console.log(`   Found Offer ${offer.id} from ${offer.supplier_name}`);
            // 4. Accept Quote (Operator)
            console.log(`4. Operator: Accepting Offer ${offer.id}...`);
            await operatorClient.post(`/quotes/${offer.id}/accept`, {});
            console.log("   Offer Accepted.");
        }

        // 5. List Assigned Trips (Partner)
        console.log("5. Partner: Listing Assigned Trips...");
        const assigned = await partnerClient.get('/trips').then(r => r.data);
        const found = assigned.find(t => t.id === tripID);

        if (found) {
            console.log(`SUCCESS: Trip ${found.booking_reference} is now in Assigned Trips with status ${found.status}`);
        } else {
            console.error(`FAILURE: Trip ID ${tripID} NOT found in Assigned Trips.`);
            console.log("DEBUG: Assigned Trips found:", JSON.stringify(assigned.map(t => ({ id: t.id, ref: t.booking_reference })), null, 2));
        }

    } catch (error) {
        console.error("Flow Failed:", error.message);
        if (error.response) console.error(error.response.data);
    }
}

runFlow();
