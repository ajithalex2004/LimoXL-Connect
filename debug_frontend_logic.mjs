import axios from 'axios';

async function testFrontendLogic() {
    try {
        console.log("Fetching RFQs...");
        const response = await axios.get('http://localhost:8080/api/partner/rfqs');
        const data = response.data;

        console.log(`Received ${data.length} items`);

        // Simulate Frontend Mapping to check for crashes
        data.forEach(trip => {
            console.log(`Processing Trip ${trip.id}`);
            const serviceType = trip.service_type?.replace('_', ' ') || 'ONE WAY';
            const pickupTime = new Date(trip.pickup_time).toLocaleDateString();
            console.log(` - Service: ${serviceType}`);
            console.log(` - Pickup: ${pickupTime}`);
        });

        console.log("Frontend logic simulation PASSED");
    } catch (error) {
        console.error("Frontend logic simulation FAILED", error.message);
    }
}

testFrontendLogic();
