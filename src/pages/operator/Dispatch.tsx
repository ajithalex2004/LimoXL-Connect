import { useState, useEffect } from 'react';
import { operatorService, type Trip } from '../../services/api';
import { api } from '../../lib/auth';
import {
    MapPin,
    User,
    Calendar,
    CheckCircle,
    Clock,
    AlertCircle,
    X,
    Users,
    Globe,
    Plus,
    Car,
    DollarSign
} from 'lucide-react';

interface Company {
    id: string;
    name: string;
    type: string;
}

plate_number: string;
type: string;
permit_expiry ?: string;
insurance_expiry ?: string;
}

interface Driver {
    id: string;
    name: string;
    license_expiry?: string;
    itc_permit_expiry?: string;
}

const Dispatch = () => {
    const [trips, setTrips] = useState<Trip[]>([]);
    const [partners, setPartners] = useState<Company[]>([]);
    const [loading, setLoading] = useState(true);

    // Assign Modal State (Outsource)
    const [showAssignModal, setShowAssignModal] = useState(false);
    const [selectedTripId, setSelectedTripId] = useState<string | null>(null);
    const [visibilityType, setVisibilityType] = useState<'BROADCAST' | 'DIRECT'>('BROADCAST');
    const [selectedPartners, setSelectedPartners] = useState<string[]>([]);
    const [submitting, setSubmitting] = useState(false);

    // Dispatch Modal State (Internal)
    const [showDispatchModal, setShowDispatchModal] = useState(false);
    const [myVehicles, setMyVehicles] = useState<Vehicle[]>([]);
    const [myDrivers, setMyDrivers] = useState<Driver[]>([]);
    const [selectedVehicle, setSelectedVehicle] = useState('');
    const [selectedDriver, setSelectedDriver] = useState('');

    // Create Trip Modal State
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [newTrip, setNewTrip] = useState({
        passenger_name: '',
        passenger_phone: '',
        pickup_time: '',
        pickup_zone: '',
        dropoff_zone: '',
        requested_vehicle_type: 'Sedan',
        price: '',
        pickup_landmark: '',
        dropoff_landmark: ''
    });

    useEffect(() => {
        loadTrips();
        loadPartners();
        loadFleet();
    }, []);

    const loadTrips = async () => {
        try {
            const data = await operatorService.listOperatorTrips();
            setTrips(data);
        } catch (error) {
            console.error("Failed to load operator trips", error);
        } finally {
            setLoading(false);
        }
    };

    const loadPartners = async () => {
        try {
            // For trip assignment, we need to use companies table IDs
            // because trip_access references companies, not outsource_companies
            const res = await api.get<any[]>('/companies');
            const supplyPartners = res.data.filter((c: any) =>
                c.type === 'SUPPLY' || c.type === 'BOTH'
            );

            // Remove duplicates by name (keep first occurrence)
            const uniquePartners = supplyPartners.filter((partner: any, index: number, self: any[]) =>
                index === self.findIndex((p: any) => p.name === partner.name)
            );

            setPartners(uniquePartners);
        } catch (err) {
            console.error("Failed to load partners", err);
        }
    };

    const loadFleet = async () => {
        try {
            const [vehicles, drivers] = await Promise.all([
                operatorService.listVehicles(),
                operatorService.listDrivers()
            ]);
            setMyVehicles(vehicles);
            setMyDrivers(drivers);
        } catch (err) {
            console.error("Failed to load fleet", err);
        }
    };

    const openAssignModal = (tripId: string) => {
        setSelectedTripId(tripId);
        setVisibilityType('BROADCAST');
        setSelectedPartners([]);
        setShowAssignModal(true);
    };

    const openDispatchModal = (tripId: string) => {
        setSelectedTripId(tripId);
        setSelectedVehicle('');
        setSelectedDriver('');
        setShowDispatchModal(true);
    };

    const handleAssignSubmit = async () => {
        if (!selectedTripId) return;

        if (visibilityType === 'DIRECT' && selectedPartners.length === 0) {
            alert("Please select at least one partner for direct assignment");
            return;
        }

        setSubmitting(true);
        try {
            const partnerIds = visibilityType === 'DIRECT' ? selectedPartners : undefined;
            await operatorService.assignOutsource(selectedTripId, partnerIds);
            await loadTrips();
            setShowAssignModal(false);
        } catch (error) {
            console.error("Failed to assign outsource", error);
            alert("Failed to assign outsource. Please try again.");
        } finally {
            setSubmitting(false);
        }
    };

    const handleDispatchSubmit = async () => {
        if (!selectedTripId || !selectedDriver || !selectedVehicle) {
            alert("Please select both a driver and a vehicle.");
            return;
        }

        setSubmitting(true);
        try {
            await operatorService.dispatchTrip(selectedTripId, selectedDriver, selectedVehicle);
            await loadTrips();
            setShowDispatchModal(false);
        } catch (error) {
            console.error("Failed to dispatch trip", error);
            alert("Failed to dispatch trip.");
        } finally {
            setSubmitting(false);
        }
    };

    const togglePartner = (id: string) => {
        if (selectedPartners.includes(id)) {
            setSelectedPartners(selectedPartners.filter(p => p !== id));
        } else {
            setSelectedPartners([...selectedPartners, id]);
        }
    };

    const handleCreateTrip = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await operatorService.createTrip({
                ...newTrip,
                price: parseFloat(newTrip.price) || 0,
                pickup_time: new Date(newTrip.pickup_time).toISOString()
            });
            setShowCreateModal(false);
            setNewTrip({
                passenger_name: '',
                passenger_phone: '',
                pickup_time: '',
                pickup_zone: '',
                dropoff_zone: '',
                requested_vehicle_type: 'Sedan',
                price: '',
                pickup_landmark: '',
                dropoff_landmark: ''
            });
            loadTrips();
        } catch (error) {
            console.error("Failed to create trip", error);
            alert("Failed to create trip");
        }
    };


    if (loading) return <div className="p-8 text-center text-gray-500">Loading Bookings...</div>;

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">All Trips</h1>
                    <p className="text-gray-500">Manage your bookings, dispatch, and farm-out requests</p>
                </div>
                <button
                    onClick={() => setShowCreateModal(true)}
                    className="flex items-center gap-2 px-4 py-2 bg-black text-white rounded-lg hover:bg-gray-800"
                >
                    <Plus className="h-4 w-4" /> New Booking
                </button>
            </div>

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
                <table className="w-full text-left text-sm">
                    <thead className="bg-gray-50 border-b border-gray-200">
                        <tr>
                            <th className="px-6 py-4 font-semibold text-gray-900">Booking Ref</th>
                            <th className="px-6 py-4 font-semibold text-gray-900">Pickup Time</th>
                            <th className="px-6 py-4 font-semibold text-gray-900">Passenger</th>
                            <th className="px-6 py-4 font-semibold text-gray-900">Route</th>
                            <th className="px-6 py-4 font-semibold text-gray-900">Status</th>
                            <th className="px-6 py-4 font-semibold text-gray-900 text-right">Action</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {trips.map((trip) => (
                            <tr key={trip.id} className="hover:bg-gray-50 transition-colors">
                                <td className="px-6 py-4 font-medium text-gray-900">
                                    {trip.booking_reference}
                                </td>
                                <td className="px-6 py-4 text-gray-600">
                                    <div className="flex items-center gap-2">
                                        <Calendar className="h-4 w-4 text-gray-400" />
                                        <span>
                                            {new Date(trip.pickup_time).toLocaleDateString()}
                                            <span className="text-gray-400 mx-1">|</span>
                                            {new Date(trip.pickup_time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-gray-600">
                                    <div className="flex items-center gap-2">
                                        <User className="h-4 w-4 text-gray-400" />
                                        <span>{trip.passenger_name}</span>
                                    </div>
                                </td>
                                <td className="px-6 py-4 text-gray-600">
                                    <div className="flex flex-col gap-1">
                                        <span className="flex items-center gap-1">
                                            <MapPin className="h-3 w-3 text-green-600" /> {trip.pickup_zone}
                                            {trip.pickup_landmark && <span className="text-gray-400 text-xs ml-1">({trip.pickup_landmark})</span>}
                                        </span>
                                        <span className="flex items-center gap-1">
                                            <MapPin className="h-3 w-3 text-red-600" /> {trip.dropoff_zone}
                                            {trip.dropoff_landmark && <span className="text-gray-400 text-xs ml-1">({trip.dropoff_landmark})</span>}
                                        </span>
                                    </div>
                                </td>
                                <td className="px-6 py-4">
                                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border ${(trip.status === 'BOOKED' || trip.status === 'CREATED') ? 'bg-blue-50 text-blue-700 border-blue-200' :
                                        trip.status === 'MARKETPLACE_SEARCH' ? 'bg-amber-50 text-amber-700 border-amber-200' :
                                            trip.status === 'OFFERED' ? 'bg-indigo-50 text-indigo-700 border-indigo-200' :
                                                (trip.status === 'ACCEPTED' || trip.status === 'OUTSOURCE_ASSIGNED' || trip.status === 'DRIVER_ASSIGNED') ? 'bg-emerald-50 text-emerald-700 border-emerald-200' :
                                                    'bg-gray-50 text-gray-600 border-gray-200'
                                        }`}>
                                        {(trip.status === 'BOOKED' || trip.status === 'CREATED') && <Clock className="h-3 w-3" />}
                                        {trip.status === 'MARKETPLACE_SEARCH' && <AlertCircle className="h-3 w-3" />}
                                        {trip.status === 'OFFERED' && <DollarSign className="h-3 w-3" />}
                                        {(trip.status === 'ACCEPTED' || trip.status === 'OUTSOURCE_ASSIGNED' || trip.status === 'DRIVER_ASSIGNED') && <CheckCircle className="h-3 w-3" />}
                                        {trip.status.replace('_', ' ')}
                                    </span>
                                </td>
                                <td className="px-6 py-4 text-right">
                                    {(trip.status === 'BOOKED' || trip.status === 'CREATED' || trip.status === 'NEW' || trip.status === 'APPROVED') && (
                                        <div className="flex justify-end gap-2">
                                            <button
                                                onClick={() => openAssignModal(trip.id)}
                                                className="inline-flex items-center gap-2 px-3 py-1.5 bg-white border border-gray-300 text-gray-700 hover:bg-gray-50 text-xs font-bold uppercase tracking-wider rounded-lg shadow-sm transition-all"
                                            >
                                                Outsource <Globe className="h-3 w-3" />
                                            </button>
                                            <button
                                                onClick={() => openDispatchModal(trip.id)}
                                                className="inline-flex items-center gap-2 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white text-xs font-bold uppercase tracking-wider rounded-lg shadow-sm transition-all active:scale-95"
                                            >
                                                Dispatch <Car className="h-3 w-3" />
                                            </button>
                                        </div>
                                    )}
                                    {trip.status === 'OFFERED' && (
                                        <div className="flex justify-end gap-2">
                                            <a
                                                href="/quotes"
                                                className="inline-flex items-center gap-2 px-3 py-1.5 bg-indigo-600 hover:bg-indigo-700 text-white text-xs font-bold uppercase tracking-wider rounded-lg shadow-sm transition-all"
                                            >
                                                View Quotes <DollarSign className="h-3 w-3" />
                                            </a>
                                        </div>
                                    )}
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>

            {/* Assign Modal (Outsource) */}
            {showAssignModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" onClick={() => setShowAssignModal(false)}>
                    <div className="bg-white rounded-xl shadow-xl max-w-lg w-full p-6" onClick={e => e.stopPropagation()}>
                        <div className="flex justify-between items-start mb-6">
                            <div>
                                <h3 className="text-lg font-bold text-gray-900">Assign Outsource</h3>
                                <p className="text-sm text-gray-500">Select how you want to farm-out this trip</p>
                            </div>
                            <button onClick={() => setShowAssignModal(false)} className="text-gray-400 hover:text-gray-500"><X className="h-5 w-5" /></button>
                        </div>
                        {/* ... Existing Outsource Modal logic ... */}
                        <div className="space-y-4 mb-8">
                            {/* Broadcast Option */}
                            <div
                                className={`p-4 border rounded-xl cursor-pointer transition-colors ${visibilityType === 'BROADCAST' ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-blue-300'}`}
                                onClick={() => setVisibilityType('BROADCAST')}
                            >
                                <div className="flex items-center gap-3">
                                    <div className={`w-5 h-5 rounded-full border flex items-center justify-center ${visibilityType === 'BROADCAST' ? 'border-blue-600' : 'border-gray-300'}`}>
                                        {visibilityType === 'BROADCAST' && <div className="w-2.5 h-2.5 rounded-full bg-blue-600" />}
                                    </div>
                                    <div>
                                        <p className="font-medium text-gray-900 flex items-center gap-2">
                                            <Globe className="h-4 w-4 text-blue-600" /> Broadcast to All
                                        </p>
                                        <p className="text-xs text-gray-500">Visible to all onboarded outsourcing partners</p>
                                    </div>
                                </div>
                            </div>

                            {/* Direct Option */}
                            <div
                                className={`p-4 border rounded-xl cursor-pointer transition-colors ${visibilityType === 'DIRECT' ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-blue-300'}`}
                                onClick={() => setVisibilityType('DIRECT')}
                            >
                                <div className="flex items-center gap-3">
                                    <div className={`w-5 h-5 rounded-full border flex items-center justify-center ${visibilityType === 'DIRECT' ? 'border-blue-600' : 'border-gray-300'}`}>
                                        {visibilityType === 'DIRECT' && <div className="w-2.5 h-2.5 rounded-full bg-blue-600" />}
                                    </div>
                                    <div>
                                        <p className="font-medium text-gray-900 flex items-center gap-2">
                                            <Users className="h-4 w-4 text-purple-600" /> Select Specific Partners
                                        </p>
                                        <p className="text-xs text-gray-500">Only selected partners will see this RFQ</p>
                                    </div>
                                </div>
                            </div>

                            {/* Partner List */}
                            {visibilityType === 'DIRECT' && (
                                <div className="bg-gray-50 p-4 rounded-xl border border-gray-200 max-h-48 overflow-y-auto">
                                    <p className="text-xs font-semibold text-gray-500 mb-3 uppercase tracking-wider">Select Partners</p>
                                    <div className="space-y-2">
                                        {partners.length === 0 ? (
                                            <p className="text-sm text-gray-500 italic">No approved partners found.</p>
                                        ) : (
                                            partners.map(partner => (
                                                <label key={partner.id} className="flex items-center gap-3 p-2 bg-white rounded-lg border border-gray-200 cursor-pointer hover:border-blue-300">
                                                    <input
                                                        type="checkbox"
                                                        checked={selectedPartners.includes(partner.id)}
                                                        onChange={() => togglePartner(partner.id)}
                                                        className="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
                                                    />
                                                    <span className="text-sm font-medium text-gray-700">{partner.name}</span>
                                                </label>
                                            ))
                                        )}
                                    </div>
                                </div>
                            )}
                        </div>
                        <button
                            onClick={handleAssignSubmit}
                            disabled={submitting}
                            className="w-full py-3 bg-blue-600 hover:bg-blue-700 text-white font-bold rounded-xl shadow-md disabled:opacity-50"
                        >
                            {submitting ? 'Processing...' : 'Confirm Assignment'}
                        </button>
                    </div>
                </div>
            )}

            {/* Dispatch Modal (Internal) */}
            {showDispatchModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm" onClick={() => setShowDispatchModal(false)}>
                    <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6" onClick={e => e.stopPropagation()}>
                        <div className="flex justify-between items-start mb-6">
                            <div>
                                <h3 className="text-lg font-bold text-gray-900">Internal Dispatch</h3>
                                <p className="text-sm text-gray-500">Assign own driver and vehicle</p>
                            </div>
                            <button onClick={() => setShowDispatchModal(false)} className="text-gray-400 hover:text-gray-500"><X className="h-5 w-5" /></button>
                        </div>

                        <div className="space-y-4 mb-6">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Vehicle</label>
                                <select
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={selectedVehicle}
                                    onChange={e => setSelectedVehicle(e.target.value)}
                                >
                                    <option value="">Select Vehicle...</option>
                                    {myVehicles.map(v => {
                                        const isExpiring = v.permit_expiry && new Date(v.permit_expiry) < new Date(Date.now() + 30 * 24 * 60 * 60 * 1000);
                                        const isExpired = v.permit_expiry && new Date(v.permit_expiry) < new Date();
                                        return (
                                            <option key={v.id} value={v.id} className={isExpired ? 'text-red-600' : isExpiring ? 'text-amber-600' : ''}>
                                                {v.plate_number} - {v.type} {isExpired ? '(EXPIRED)' : isExpiring ? '(EXPIRING SOON)' : ''}
                                            </option>
                                        );
                                    })}
                                </select>
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Driver</label>
                                <select
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={selectedDriver}
                                    onChange={e => setSelectedDriver(e.target.value)}
                                >
                                    <option value="">Select Driver...</option>
                                    {myDrivers.map(d => {
                                        const isExpiring = d.license_expiry && new Date(d.license_expiry) < new Date(Date.now() + 30 * 24 * 60 * 60 * 1000);
                                        const isExpired = d.license_expiry && new Date(d.license_expiry) < new Date();
                                        return (
                                            <option key={d.id} value={d.id} className={isExpired ? 'text-red-600' : isExpiring ? 'text-amber-600' : ''}>
                                                {d.name} {isExpired ? '(EXPIRED)' : isExpiring ? '(EXPIRING SOON)' : ''}
                                            </option>
                                        );
                                    })}
                                </select>
                            </div>
                        </div>

                        <button
                            onClick={handleDispatchSubmit}
                            disabled={submitting}
                            className="w-full py-3 bg-black hover:bg-gray-800 text-white font-bold rounded-xl shadow-md disabled:opacity-50"
                        >
                            {submitting ? 'Dispatching...' : 'Confirm Dispatch'}
                        </button>
                    </div>
                </div>
            )}

            {/* Create Trip Modal */}
            {showCreateModal && (
                <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
                    <div className="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
                        <div className="flex justify-between items-start mb-6">
                            <h2 className="text-xl font-bold">New Booking</h2>
                            <button onClick={() => setShowCreateModal(false)}><X className="h-5 w-5 text-gray-500" /></button>
                        </div>
                        <form onSubmit={handleCreateTrip} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Passenger Name</label>
                                <input
                                    type="text"
                                    required
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={newTrip.passenger_name}
                                    onChange={e => setNewTrip({ ...newTrip, passenger_name: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Passenger Phone</label>
                                <input
                                    type="tel"
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={newTrip.passenger_phone}
                                    onChange={e => setNewTrip({ ...newTrip, passenger_phone: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Pickup Time</label>
                                <input
                                    type="datetime-local"
                                    required
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={newTrip.pickup_time}
                                    onChange={e => setNewTrip({ ...newTrip, pickup_time: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Pickup Address / Zone</label>
                                <input
                                    type="text"
                                    required
                                    placeholder="e.g. Dubai Mall"
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={newTrip.pickup_zone}
                                    onChange={e => setNewTrip({ ...newTrip, pickup_zone: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Pickup Landmark (Internal)</label>
                                <input
                                    type="text"
                                    placeholder="e.g. Entrance 4, Level LG"
                                    className="w-full px-3 py-2 border rounded-lg bg-blue-50/30 font-medium"
                                    value={newTrip.pickup_landmark}
                                    onChange={e => setNewTrip({ ...newTrip, pickup_landmark: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Dropoff Address / Zone</label>
                                <input
                                    type="text"
                                    required
                                    placeholder="e.g. Marina 101"
                                    className="w-full px-3 py-2 border rounded-lg"
                                    value={newTrip.dropoff_zone}
                                    onChange={e => setNewTrip({ ...newTrip, dropoff_zone: e.target.value })}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Dropoff Landmark (Internal)</label>
                                <input
                                    type="text"
                                    placeholder="e.g. Main Lobby Gate"
                                    className="w-full px-3 py-2 border rounded-lg bg-blue-50/30 font-medium"
                                    value={newTrip.dropoff_landmark}
                                    onChange={e => setNewTrip({ ...newTrip, dropoff_landmark: e.target.value })}
                                />
                            </div>
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Vehicle Type</label>
                                    <select
                                        className="w-full px-3 py-2 border rounded-lg"
                                        value={newTrip.requested_vehicle_type}
                                        onChange={e => setNewTrip({ ...newTrip, requested_vehicle_type: e.target.value })}
                                    >
                                        <option value="Sedan">Sedan</option>
                                        <option value="SUV">SUV</option>
                                        <option value="Van">Van</option>
                                        <option value="Luxury">Luxury</option>
                                    </select>
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-gray-700 mb-1">Expected Price (AED)</label>
                                    <input
                                        type="number"
                                        className="w-full px-3 py-2 border rounded-lg"
                                        value={newTrip.price}
                                        onChange={e => setNewTrip({ ...newTrip, price: e.target.value })}
                                    />
                                </div>
                            </div>

                            <button type="submit" className="w-full py-3 bg-black text-white font-bold rounded-lg mt-2">
                                Create Booking
                            </button>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Dispatch;
