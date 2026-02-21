import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { MapPin, Phone, User, Navigation, CheckCircle } from 'lucide-react';
import axios from 'axios';

// Simplify for standalone usage without full api service if needed, 
// but we can reuse the existing service if we export a method or just use axios directly.
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

interface TripDetails {
    id: string;
    booking_reference: string;
    pickup_zone: string;
    dropoff_zone: string;
    pickup_time: string;
    passenger_name: string;
    passenger_phone: string;
    status: string;
}

const SecureDriverLink = () => {
    const { token } = useParams<{ token: string }>();
    const [trip, setTrip] = useState<TripDetails | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        if (token) {
            fetchTripStatus();
        }
    }, [token]);

    const fetchTripStatus = async () => {
        try {
            const res = await axios.get(`${API_URL}/status/${token}`);
            setTrip(res.data);
        } catch (err) {
            console.error("Failed to load trip", err);
            setError("Invalid or expired link");
        } finally {
            setLoading(false);
        }
    };

    const updateStatus = async (newStatus: string) => {
        if (!token) return;
        try {
            // Optimistic update
            if (trip) setTrip({ ...trip, status: newStatus });

            await axios.post(`${API_URL}/status/${token}`, { status: newStatus });
            fetchTripStatus(); // Refresh to ensure sync
        } catch (err) {
            console.error("Failed to update status", err);
            alert("Failed to update status");
        }
    };

    if (loading) return <div className="flex items-center justify-center min-h-screen bg-gray-100">Loading...</div>;
    if (error) return <div className="flex items-center justify-center min-h-screen bg-gray-100 text-red-600 font-bold">{error}</div>;
    if (!trip) return null;

    return (
        <div className="min-h-screen bg-gray-50 flex flex-col">
            {/* Header */}
            <div className="bg-white shadow-sm p-4 sticky top-0 z-10">
                <h1 className="text-xl font-bold text-gray-900 text-center">Limo Driver App</h1>
                <div className="mt-2 flex justify-center">
                    <span className="bg-emerald-100 text-emerald-800 text-xs font-bold px-3 py-1 rounded-full uppercase tracking-wider">
                        {trip.status.replace('_', ' ')}
                    </span>
                </div>
            </div>

            <div className="flex-1 p-4 space-y-6 max-w-lg mx-auto w-full">
                {/* Trip Info Card */}
                <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
                    <div className="p-4 border-b border-gray-100 bg-gray-50/50">
                        <p className="text-sm text-gray-500">Booking Reference</p>
                        <p className="text-lg font-bold text-gray-900">{trip.booking_reference}</p>
                    </div>

                    <div className="p-4 space-y-4">
                        <div className="flex items-start gap-3">
                            <MapPin className="h-5 w-5 text-green-600 mt-0.5 shrink-0" />
                            <div>
                                <p className="text-sm font-medium text-gray-900">Pickup</p>
                                <p className="text-gray-600">{trip.pickup_zone}</p>
                                <p className="text-xs text-gray-400 mt-1">
                                    {new Date(trip.pickup_time).toLocaleString()}
                                </p>
                            </div>
                        </div>

                        <div className="flex items-start gap-3">
                            <MapPin className="h-5 w-5 text-red-600 mt-0.5 shrink-0" />
                            <div>
                                <p className="text-sm font-medium text-gray-900">Dropoff</p>
                                <p className="text-gray-600">{trip.dropoff_zone}</p>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Passenger Info */}
                <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-4">
                    <div className="flex items-center gap-3 mb-3">
                        <User className="h-5 w-5 text-gray-400" />
                        <div>
                            <p className="text-sm text-gray-500">Passenger</p>
                            <p className="font-medium text-gray-900">{trip.passenger_name}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-3">
                        <Phone className="h-5 w-5 text-gray-400" />
                        <div>
                            <p className="text-sm text-gray-500">Contact</p>
                            <a href={trip.passenger_phone ? `tel:${trip.passenger_phone}` : '#'} className={`font-medium ${trip.passenger_phone ? 'text-blue-600 hover:underline' : 'text-gray-400 cursor-not-allowed'}`}>
                                {trip.passenger_phone || 'N/A'}
                            </a>
                        </div>
                    </div>
                </div>

                {/* Actions */}
                <div className="space-y-3 pb-8">
                    {/* Allow starting if Driver Assigned OR (Edge case) just Assigned/Marketplace but has link */}
                    {['DRIVER_ASSIGNED', 'OUTSOURCE_ASSIGNED', 'MARKETPLACE_SEARCH', 'ACCEPTED'].includes(trip.status) && (
                        <button
                            onClick={() => updateStatus('EN_ROUTE')}
                            className="w-full py-4 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded-xl shadow-md flex items-center justify-center gap-2 active:scale-95 transition-transform"
                        >
                            <Navigation className="h-6 w-6" />
                            Start / En Route to Pickup
                        </button>
                    )}

                    {trip.status === 'EN_ROUTE' && (
                        <button
                            onClick={() => updateStatus('IN_TRIP')}
                            className="w-full py-4 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded-xl shadow-md flex items-center justify-center gap-2 active:scale-95 transition-transform"
                        >
                            <User className="h-6 w-6" />
                            Passenger On Board
                        </button>
                    )}

                    {trip.status === 'IN_TRIP' && (
                        <button
                            onClick={() => updateStatus('COMPLETED')}
                            className="w-full py-4 bg-green-600 hover:bg-green-700 text-white font-bold rounded-xl shadow-md flex items-center justify-center gap-2 active:scale-95 transition-transform"
                        >
                            <CheckCircle className="h-6 w-6" />
                            Complete Trip
                        </button>
                    )}

                    {trip.status === 'COMPLETED' && (
                        <div className="bg-green-50 text-green-800 p-4 rounded-xl text-center font-bold border border-green-200">
                            Trip Completed Successfully
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default SecureDriverLink;
