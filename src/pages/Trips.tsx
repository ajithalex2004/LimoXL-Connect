import { useState, useEffect } from 'react';
import { partnerService, type Trip, type Driver, type Vehicle } from '../services/api';
import {
    MapPin,
    User,
    Phone,
    Share2,
    Copy,
    MessageSquare,
    MessageCircle, // WhatsApp
    Check,
    X,
    Car
} from 'lucide-react';

const Trips = () => {
    const [trips, setTrips] = useState<Trip[]>([]);
    const [drivers, setDrivers] = useState<Driver[]>([]);
    const [vehicles, setVehicles] = useState<Vehicle[]>([]);
    const [loading, setLoading] = useState(true);
    const [selectedTrip, setSelectedTrip] = useState<Trip | null>(null);
    const [showLinkModal, setShowLinkModal] = useState(false);
    const [showAssignModal, setShowAssignModal] = useState(false);
    const [assignForm, setAssignForm] = useState({
        driverName: '',
        driverPhone: '',
        vehiclePlate: '',
        vehicleModel: ''
    });
    const [copied, setCopied] = useState(false);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const [tripsData, driversData, vehiclesData] = await Promise.all([
                partnerService.listAssignedTrips().catch(() => []),
                partnerService.listDrivers().catch(() => []),
                partnerService.listVehicles().catch(() => [])
            ]);
            setTrips(tripsData || []);
            setDrivers(driversData || []);
            setVehicles(vehiclesData || []);
        } catch (error) {
            console.error("Failed to load data", error);
            setTrips([]);
            setDrivers([]);
            setVehicles([]);
        } finally {
            setLoading(false);
        }
    };

    const handleShareClick = (trip: Trip) => {
        setSelectedTrip(trip);
        setShowLinkModal(true);
        setCopied(false);
    };

    const handleAssignClick = (trip: Trip) => {
        setSelectedTrip(trip);
        setAssignForm({
            driverName: '',
            driverPhone: '',
            vehiclePlate: '',
            vehicleModel: ''
        });
        setShowAssignModal(true);
    };

    const handleAssignSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!selectedTrip) return;

        try {
            await partnerService.assignDriver({
                trip_id: selectedTrip.id,
                driver_name: assignForm.driverName,
                driver_phone: assignForm.driverPhone,
                vehicle_model: assignForm.vehicleModel,
                vehicle_plate: assignForm.vehiclePlate
            });

            // Reload all data to get updated trips with token
            const [updatedTrips, driversData, vehiclesData] = await Promise.all([
                partnerService.listAssignedTrips(),
                partnerService.listDrivers(),
                partnerService.listVehicles()
            ]);
            setTrips(updatedTrips);
            setDrivers(driversData);
            setVehicles(vehiclesData);

            // Find the updated trip to get the new token
            const updatedTrip = updatedTrips.find(t => t.id === selectedTrip.id);
            if (updatedTrip && updatedTrip.driver_link_token) {
                setSelectedTrip(updatedTrip);
                setShowAssignModal(false);
                // Automatically show the Share Link modal
                setShowLinkModal(true);
            } else {
                setShowAssignModal(false);
                alert('Driver assigned successfully!');
            }
        } catch (error) {
            console.error("Failed to assign driver", error);
            alert("Failed to assign driver. Please try again.");
        }
    };

    const getDriverLink = (trip: Trip) => {
        // Use the current origin so it works if accessed via IP or domain
        // Ensure we only use the real token. If missing, it will be empty (handled by UI showing Assign button)
        const token = trip.driver_link_token || '';
        return `${window.location.origin}/link/${token}`;
    };

    const copyToClipboard = () => {
        if (!selectedTrip) return;
        const link = getDriverLink(selectedTrip);
        navigator.clipboard.writeText(link);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const sendSMS = () => {
        if (!selectedTrip) return;
        const link = getDriverLink(selectedTrip);
        const body = `Trip Assignment: ${selectedTrip.booking_reference}. Link: ${link}`;
        window.open(`sms:?&body=${encodeURIComponent(body)}`, '_self');
    };

    const sendWhatsApp = () => {
        if (!selectedTrip) return;
        const link = getDriverLink(selectedTrip);
        const text = `Trip Assignment: ${selectedTrip.booking_reference}. Link: ${link}`;
        window.open(`https://wa.me/?text=${encodeURIComponent(text)}`, '_blank');
    };

    if (loading) return <div className="p-8 text-center text-gray-500">Loading Assigned Trips...</div>;

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-gray-900">Assigned Trips</h1>
                <p className="text-gray-500">Manage ongoing trips and driver assignments</p>
            </div>

            <div className="grid gap-4">
                {trips.map((trip) => (
                    <div key={trip.id} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 flex flex-col lg:flex-row justify-between gap-6 hover:border-emerald-200 transition-colors">
                        <div className="flex-1 space-y-4">
                            <div className="flex items-center gap-3">
                                <span className="bg-emerald-100 text-emerald-800 text-xs font-semibold px-2.5 py-0.5 rounded">
                                    {trip.booking_reference}
                                </span>
                                <span className={`text-xs font-bold px-2 py-0.5 rounded border ${trip.status === 'DRIVER_ASSIGNED' ? 'bg-blue-50 text-blue-600 border-blue-100' :
                                    trip.status === 'EN_ROUTE' ? 'bg-amber-50 text-amber-600 border-amber-100' :
                                        trip.status === 'IN_TRIP' ? 'bg-purple-50 text-purple-600 border-purple-100' :
                                            'bg-gray-50 text-gray-600 border-gray-100'
                                    }`}>
                                    {trip.status.replace('_', ' ')}
                                </span>
                            </div>

                            <div className="grid md:grid-cols-2 gap-6">
                                <div className="space-y-2">
                                    <div className="flex items-start gap-2">
                                        <MapPin className="h-4 w-4 text-green-600 mt-1 shrink-0" />
                                        <div>
                                            <p className="text-sm font-medium text-gray-900">
                                                {trip.pickup_zone}
                                                {trip.pickup_landmark && <span className="text-gray-400 font-normal ml-1">({trip.pickup_landmark})</span>}
                                            </p>
                                            <p className="text-xs text-gray-500">
                                                {new Date(trip.pickup_time).toLocaleDateString()} {new Date(trip.pickup_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                            </p>
                                        </div>
                                    </div>
                                    <div className="flex items-start gap-2">
                                        <MapPin className="h-4 w-4 text-red-600 mt-1 shrink-0" />
                                        <div>
                                            <p className="text-sm font-medium text-gray-900">
                                                {trip.dropoff_zone}
                                                {trip.dropoff_landmark && <span className="text-gray-400 font-normal ml-1">({trip.dropoff_landmark})</span>}
                                            </p>
                                        </div>
                                    </div>
                                </div>

                                <div className="space-y-2">
                                    <div className="flex items-center gap-2">
                                        <User className="h-4 w-4 text-gray-400" />
                                        <span className="text-sm text-gray-600">{trip.passenger_name}</span>
                                    </div>
                                    <div className="flex items-center gap-2">
                                        <Phone className="h-4 w-4 text-gray-400" />
                                        <span className="text-sm text-gray-600">{trip.passenger_phone}</span>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="flex flex-col justify-center items-end border-t lg:border-t-0 lg:border-l border-gray-100 lg:pl-6 pt-4 lg:pt-0">
                            {/* Show Assign if new OR if token is missing (migration fix) */}
                            {trip.status === 'OUTSOURCE_ASSIGNED' || !trip.driver_link_token ? (
                                <button
                                    onClick={() => handleAssignClick(trip)}
                                    className="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-lg transition-colors flex items-center gap-2 shadow-sm w-full lg:w-auto justify-center"
                                >
                                    <Car className="h-4 w-4" />
                                    {trip.driver_link_token ? 'Re-Assign Driver' : 'Assign Driver'}
                                </button>
                            ) : (
                                <button
                                    onClick={() => handleShareClick(trip)}
                                    className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium py-2 px-4 rounded-lg transition-colors flex items-center gap-2 shadow-sm w-full lg:w-auto justify-center"
                                >
                                    <Share2 className="h-4 w-4" />
                                    Share Driver Link
                                </button>
                            )}
                        </div>
                    </div>
                ))}

                {
                    trips.length === 0 && (
                        <div className="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
                            <p className="text-gray-500">No assigned trips found.</p>
                        </div>
                    )
                }
            </div >

            {/* Assign Driver Modal */}
            {
                showAssignModal && selectedTrip && (
                    <div
                        className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm transition-opacity"
                        onClick={(e) => {
                            if (e.target === e.currentTarget) setShowAssignModal(false);
                        }}
                    >
                        <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6 animate-in fade-in zoom-in duration-200">
                            <div className="flex justify-between items-start mb-4">
                                <div>
                                    <h3 className="text-lg font-bold text-gray-900">Assign Driver & Vehicle</h3>
                                    <p className="text-sm text-gray-500">Trip: {selectedTrip.booking_reference}</p>
                                </div>
                                <button
                                    onClick={() => setShowAssignModal(false)}
                                    className="text-gray-400 hover:text-gray-500 p-2 hover:bg-gray-100 rounded-full transition-colors"
                                    aria-label="Close modal"
                                >
                                    <X className="h-5 w-5" />
                                </button>
                            </div>

                            <div className="mb-4 space-y-3">
                                <div>
                                    <label className="block text-xs font-medium text-gray-500 mb-1">Quick Select Driver</label>
                                    <select
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                                        onChange={(e) => {
                                            const driver = drivers.find(d => d.id === e.target.value);
                                            if (driver) {
                                                const updates: any = {
                                                    driverName: driver.name,
                                                    driverPhone: driver.phone
                                                };

                                                // Auto-select vehicle if assigned
                                                if (driver.current_vehicle_id) {
                                                    const vehicle = vehicles.find(v => v.id === driver.current_vehicle_id);
                                                    if (vehicle) {
                                                        updates.vehicleModel = vehicle.model;
                                                        updates.vehiclePlate = vehicle.plate_number;
                                                    }
                                                }
                                                setAssignForm(prev => ({ ...prev, ...updates }));
                                            }
                                        }}
                                        defaultValue=""
                                    >
                                        <option value="" disabled>Select a driver...</option>
                                        {drivers.map(d => (
                                            <option key={d.id} value={d.id}>{d.name}</option>
                                        ))}
                                    </select>
                                </div>

                                <div>
                                    <label className="block text-xs font-medium text-gray-500 mb-1">Quick Select Vehicle</label>
                                    <select
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm"
                                        onChange={(e) => {
                                            const vehicle = vehicles.find(v => v.id === e.target.value);
                                            if (vehicle) {
                                                setAssignForm(prev => ({
                                                    ...prev,
                                                    vehicleModel: vehicle.model,
                                                    vehiclePlate: vehicle.plate_number
                                                }));
                                            }
                                        }}
                                        defaultValue=""
                                    >
                                        <option value="" disabled>Select a vehicle...</option>
                                        {vehicles.map(v => (
                                            <option key={v.id} value={v.id}>{v.model} - {v.plate_number}</option>
                                        ))}
                                    </select>
                                </div>
                                <div className="border-t border-gray-100 my-4"></div>
                            </div>

                            <form onSubmit={handleAssignSubmit} className="space-y-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Driver Name</label>
                                    <input
                                        type="text"
                                        required
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-gray-50"
                                        value={assignForm.driverName}
                                        onChange={e => setAssignForm({ ...assignForm, driverName: e.target.value })}
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Driver Phone</label>
                                    <input
                                        type="tel"
                                        required
                                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-gray-50"
                                        value={assignForm.driverPhone}
                                        onChange={e => setAssignForm({ ...assignForm, driverPhone: e.target.value })}
                                    />
                                </div>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Vehicle Model</label>
                                        <input
                                            type="text"
                                            required
                                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-gray-50"
                                            value={assignForm.vehicleModel}
                                            onChange={e => setAssignForm({ ...assignForm, vehicleModel: e.target.value })}
                                        />
                                    </div>
                                    <div>
                                        <label className="block text-sm font-medium text-gray-700 mb-1">Plate Number</label>
                                        <input
                                            type="text"
                                            required
                                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500 focus:border-emerald-500 bg-gray-50"
                                            value={assignForm.vehiclePlate}
                                            onChange={e => setAssignForm({ ...assignForm, vehiclePlate: e.target.value })}
                                        />
                                    </div>
                                </div>

                                <div className="pt-2">
                                    <button
                                        type="submit"
                                        className="w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-2.5 px-4 rounded-lg transition-colors shadow-sm"
                                    >
                                        Confirm Assignment
                                    </button>
                                </div>
                            </form>
                        </div>
                    </div>
                )
            }

            {/* Share Modal */}
            {
                showLinkModal && selectedTrip && (
                    <div
                        className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm transition-opacity"
                        onClick={(e) => {
                            if (e.target === e.currentTarget) setShowLinkModal(false);
                        }}
                    >
                        <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6 animate-in fade-in zoom-in duration-200">
                            <div className="flex justify-between items-start mb-4">
                                <div>
                                    <h3 className="text-lg font-bold text-gray-900">Share Driver Link</h3>
                                    <p className="text-sm text-gray-500">Trip: {selectedTrip.booking_reference}</p>
                                </div>
                                <button
                                    onClick={() => setShowLinkModal(false)}
                                    className="text-gray-400 hover:text-gray-500 p-2 hover:bg-gray-100 rounded-full transition-colors"
                                    aria-label="Close modal"
                                >
                                    <X className="h-5 w-5" />
                                </button>
                            </div>

                            <div className="p-3 bg-gray-50 rounded-lg border border-gray-200 mb-6 break-all">
                                <p className="text-xs font-mono text-gray-600 select-all">
                                    {getDriverLink(selectedTrip)}
                                </p>
                            </div>

                            <div className="grid grid-cols-1 gap-3">
                                <button
                                    onClick={copyToClipboard}
                                    className="flex items-center justify-center gap-2 w-full py-2.5 px-4 bg-gray-900 hover:bg-gray-800 text-white rounded-lg font-medium transition-colors"
                                >
                                    {copied ? <Check className="h-5 w-5" /> : <Copy className="h-5 w-5" />}
                                    {copied ? 'Copied!' : 'Copy Link'}
                                </button>

                                <div className="grid grid-cols-2 gap-3">
                                    <button
                                        onClick={sendWhatsApp}
                                        className="flex items-center justify-center gap-2 py-2.5 px-4 bg-[#25D366] hover:bg-[#20bd5a] text-white rounded-lg font-medium transition-colors"
                                    >
                                        <MessageCircle className="h-5 w-5" />
                                        WhatsApp
                                    </button>
                                    <button
                                        onClick={sendSMS}
                                        className="flex items-center justify-center gap-2 py-2.5 px-4 bg-blue-500 hover:bg-blue-600 text-white rounded-lg font-medium transition-colors"
                                    >
                                        <MessageSquare className="h-5 w-5" />
                                        SMS
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>
                )
            }
        </div >
    );
};

export default Trips;
