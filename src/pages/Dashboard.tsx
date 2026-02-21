import { FileText, CheckCircle, Clock, AlertCircle, X, Check, Building, Users } from 'lucide-react';
import { useState, useEffect } from 'react';

import { partnerService, operatorService, type Trip } from '../services/api';
import { useAuthStore } from '../lib/auth';

const Dashboard = () => {
    const { user } = useAuthStore();
    const [activities, setActivities] = useState<Trip[]>([]);
    const [operatorTrips, setOperatorTrips] = useState<Trip[]>([]);
    const [loading, setLoading] = useState(true);
    const [quotingTrip, setQuotingTrip] = useState<Trip | null>(null);
    const [quotePrice, setQuotePrice] = useState('');
    const [quoteNotes, setQuoteNotes] = useState('');

    const isOperator = user?.role === 'ADMIN' || user?.role === 'OPS' || user?.role === 'DISPATCHER';

    useEffect(() => {
        if (!isOperator) {
            loadActivities();
        } else {
            loadOperatorData();
        }
    }, [isOperator]);

    const loadActivities = async () => {
        try {
            const [rfqs, assignedTrips] = await Promise.all([
                partnerService.listRFQs(),
                partnerService.listAssignedTrips()
            ]);

            // Ensure arrays are not null (API determines null vs empty array behavior)
            const safeRfqs = rfqs || [];
            const safeAssignedTrips = assignedTrips || [];

            // Filter out already accepted RFQs from the RFQ list (client-side specific logic)
            const acceptedIds = JSON.parse(localStorage.getItem('accepted_activities') || '[]');
            const newRfqs = safeRfqs.filter(t => !acceptedIds.includes(t.id));

            // Merge: RFQs first, then Assigned Trips, but remove duplicates by ID
            const seenIds = new Set<string>();
            const allActivities = [...newRfqs, ...safeAssignedTrips]
                .filter(trip => {
                    if (seenIds.has(trip.id)) {
                        return false; // Skip duplicate
                    }
                    seenIds.add(trip.id);
                    return true;
                })
                .sort((a, b) =>
                    new Date(b.created_at || b.pickup_time).getTime() - new Date(a.created_at || a.pickup_time).getTime()
                );

            // Filter for Dashboard "Recent RFQ Requests" (only New RFQs)
            const recentRfqs = allActivities.filter(a => a.status === 'MARKETPLACE_SEARCH');

            setActivities(recentRfqs);
        } catch (error) {
            console.error("Failed to load activities", error);
        } finally {
            setLoading(false);
        }
    };

    const loadOperatorData = async () => {
        try {
            const trips = await operatorService.listOperatorTrips();
            setOperatorTrips(trips || []);
        } catch (error) {
            console.error("Failed to load operator trips", error);
        } finally {
            setLoading(false);
        }
    };

    const handleQuoteClick = (trip: Trip) => {
        setQuotingTrip(trip);
        setQuotePrice('');
        setQuoteNotes('');
    };

    const handleSubmitQuote = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!quotingTrip) return;

        try {
            await partnerService.submitQuote({
                trip_id: quotingTrip.id,
                price: parseFloat(quotePrice),
                notes: quoteNotes
            });

            // Optimistic Update: Mark as OFFERED
            setActivities(prev => prev.map(t =>
                t.id === quotingTrip.id ? { ...t, status: 'OFFERED' } : t
            ));

            setQuotingTrip(null);
            // Background refresh
            loadActivities();
        } catch (error) {
            console.error("Failed to submit quote", error);
            alert("Failed to submit quote");
        }
    };

    const handleReject = async (id: string) => {
        // Optimistic update
        setActivities(prev => prev.filter(item => item.id !== id));
        try {
            await partnerService.rejectRFQ(id);
            // Background refresh
            loadActivities();
        } catch (error) {
            console.error("Failed to reject RFQ", error);
            // Revert on failure
            loadActivities();
            alert("Failed to reject RFQ");
        }
    };

    if (loading) return <div>Loading...</div>;

    // Operator Stats
    const operatorStats = [
        { label: 'Outsource Companies', value: '1', icon: Building, color: 'text-purple-600', bg: 'bg-purple-100' },
        { label: 'Active Drivers', value: operatorTrips.filter(t => ['DRIVER_ASSIGNED', 'EN_ROUTE', 'IN_TRIP'].includes(t.status)).length.toString(), icon: Users, color: 'text-blue-600', bg: 'bg-blue-100' },
        { label: 'Pending RFQs', value: operatorTrips.filter(t => t.status === 'MARKETPLACE_SEARCH' || t.status === 'CREATED').length.toString(), icon: Clock, color: 'text-amber-600', bg: 'bg-amber-100' },
    ];

    // Partner Stats
    const partnerStats = [
        { label: 'New RFQs', value: activities.filter(a => a.status === 'MARKETPLACE_SEARCH').length.toString(), icon: FileText, color: 'text-blue-600', bg: 'bg-blue-100' },
        { label: 'Pending Quotes', value: '1', icon: Clock, color: 'text-amber-600', bg: 'bg-amber-100' },
        { label: 'Assigned Trips', value: activities.filter(a => a.status !== 'MARKETPLACE_SEARCH').length.toString(), icon: CheckCircle, color: 'text-emerald-600', bg: 'bg-emerald-100' },
        { label: 'SLA Alerts', value: '0', icon: AlertCircle, color: 'text-red-600', bg: 'bg-red-100' },
    ];

    const stats = isOperator ? operatorStats : partnerStats;

    return (
        <div className="space-y-6">
            <div>
                <h1 className="text-2xl font-bold text-gray-900">Dashboard</h1>
                <p className="text-gray-500">Welcome to Limo XL Connect {isOperator ? '(Operator Portal)' : '(Partner Portal)'}</p>
            </div>

            <div className={`grid grid-cols-1 md:grid-cols-2 lg:grid-cols-${isOperator ? '3' : '4'} gap-4`}>
                {stats.map((stat) => (
                    <div key={stat.label} className="bg-white p-6 rounded-lg shadow-sm border border-gray-100 flex items-center gap-4">
                        <div className={`p-3 rounded-full ${stat.bg}`}>
                            <stat.icon className={`h-6 w-6 ${stat.color}`} />
                        </div>
                        <div>
                            <p className="text-sm text-gray-500 font-medium">{stat.label}</p>
                            <h3 className="text-2xl font-bold text-gray-900">{stat.value}</h3>
                        </div>
                    </div>
                ))}
            </div>

            {!isOperator && (
                <div className="bg-white rounded-lg shadow-sm border border-gray-100 p-6">
                    <h2 className="text-lg font-semibold text-gray-900 mb-4">Recent RFQ Requests</h2>
                    <div className="space-y-4">
                        {activities.length === 0 ? (
                            <p className="text-gray-500 text-center py-4">No recent activity.</p>
                        ) : (
                            activities.map((item) => (
                                <div key={item.id} className="flex flex-col md:flex-row md:items-center justify-between py-4 border-b border-gray-50 last:border-0 gap-4">
                                    <div className="flex items-start gap-3 flex-1">
                                        <div className={`mt-1.5 h-2 w-2 rounded-full shrink-0 ${item.status === 'MARKETPLACE_SEARCH' ? 'bg-blue-500' :
                                            item.status === 'OFFERED' ? 'bg-amber-500' : 'bg-emerald-500'
                                            }`}></div>

                                        <div className="space-y-1 w-full">
                                            <div className="flex justify-between items-start w-full">
                                                <div>
                                                    <p className="text-sm font-bold text-gray-900">
                                                        {item.status === 'MARKETPLACE_SEARCH' ? 'New RFQ' :
                                                            item.status === 'OFFERED' ? 'Quote Sent' : 'Assigned Trip'}: {item.booking_reference}
                                                    </p>
                                                    <p className="text-xs text-blue-600 font-medium">From {item.supplier_name || 'System'}</p>
                                                </div>
                                                <span className="text-xs font-semibold bg-gray-100 text-gray-600 px-2 py-0.5 rounded">
                                                    {item.requested_vehicle_class} {item.requested_vehicle_group}
                                                </span>
                                            </div>

                                            <div className="grid grid-cols-2 gap-x-8 gap-y-1 text-xs text-gray-500 pt-1">
                                                <div>
                                                    <span className="font-medium text-gray-700">Pickup:</span> {item.pickup_zone}
                                                    <div className="ml-1 text-gray-400">
                                                        {new Date(item.pickup_time).toLocaleDateString()}
                                                    </div>
                                                </div>
                                                <div>
                                                    <span className="font-medium text-gray-700">Dropoff:</span> {item.dropoff_zone}
                                                    <div className="ml-1 text-green-600 font-medium mt-1">
                                                        {item.status.replace('_', ' ')}
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>

                                    <div className="flex items-center gap-2 self-end md:self-center shrink-0">
                                        {item.status === 'MARKETPLACE_SEARCH' ? (
                                            <>
                                                <button
                                                    onClick={() => handleQuoteClick(item)}
                                                    className="flex items-center gap-1 text-xs bg-emerald-50 text-emerald-600 px-3 py-1.5 rounded-md hover:bg-emerald-100 transition-colors font-medium border border-emerald-200"
                                                >
                                                    <FileText className="h-3 w-3" />
                                                    Submit Quote
                                                </button>
                                                <button
                                                    onClick={() => handleReject(item.id)}
                                                    className="flex items-center gap-1 text-xs bg-red-50 text-red-600 px-3 py-1.5 rounded-md hover:bg-red-100 transition-colors font-medium border border-red-200"
                                                >
                                                    <X className="h-3 w-3" />
                                                    Reject
                                                </button>
                                            </>
                                        ) : (
                                            <span className="text-xs text-gray-500 italic px-2">
                                                {item.status === 'OFFERED' ? 'Pending Quote' : 'Accepted'}
                                            </span>
                                        )}
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>
            )}
            {isOperator && (
                <div className="bg-white rounded-lg shadow-sm border border-gray-100 p-6">
                    <h2 className="text-lg font-semibold text-gray-900 mb-4">Active Trip Status</h2>
                    <div className="overflow-x-auto">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead>
                                <tr>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Ref #</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Route</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Passenger</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
                                    <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Last Updated</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-200">
                                {operatorTrips.length === 0 ? (
                                    <tr>
                                        <td colSpan={5} className="px-4 py-4 text-center text-sm text-gray-500">No active trips found.</td>
                                    </tr>
                                ) : (
                                    operatorTrips.map(trip => (
                                        <tr key={trip.id}>
                                            <td className="px-4 py-2 text-sm font-bold text-gray-900">{trip.booking_reference}</td>
                                            <td className="px-4 py-2 text-sm text-gray-600">
                                                <div>{trip.pickup_zone}</div>
                                                <div className="text-gray-400">to {trip.dropoff_zone}</div>
                                            </td>
                                            <td className="px-4 py-2 text-sm text-gray-600">{trip.passenger_name}</td>
                                            <td className="px-4 py-2">
                                                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${trip.status === 'COMPLETED' ? 'bg-green-100 text-green-800' :
                                                    trip.status === 'EN_ROUTE' ? 'bg-amber-100 text-amber-800' :
                                                        trip.status === 'IN_TRIP' ? 'bg-purple-100 text-purple-800' :
                                                            trip.status === 'DRIVER_ASSIGNED' ? 'bg-blue-100 text-blue-800' :
                                                                'bg-gray-100 text-gray-800'
                                                    }`}>
                                                    {trip.status.replace('_', ' ')}
                                                </span>
                                            </td>
                                            <td className="px-4 py-2 text-sm text-gray-500">
                                                {trip.updated_at ? new Date(trip.updated_at).toLocaleString() : '-'}
                                                {/* Fallback to fetch current check if needed, but updated_at on trip should solve it */}
                                            </td>
                                        </tr>
                                    ))
                                )}
                            </tbody>
                        </table>
                    </div>
                </div>
            )}

            {/* Quote Modal */}
            {quotingTrip && (
                <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-50">
                    <div className="bg-white rounded-lg p-6 w-full max-w-md shadow-xl">
                        <div className="flex justify-between items-center mb-4">
                            <h3 className="text-lg font-bold">Submit Quote</h3>
                            <button onClick={() => setQuotingTrip(null)}><X className="h-5 w-5 text-gray-500" /></button>
                        </div>
                        <p className="text-sm text-gray-600 mb-4">
                            Ref: <span className="font-semibold">{quotingTrip.booking_reference}</span>
                        </p>

                        <div className="bg-gray-50 p-4 rounded-md mb-4 text-sm space-y-3 border border-gray-100">
                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <span className="text-gray-500 block text-xs uppercase tracking-wide">Pickup</span>
                                    <span className="font-medium text-gray-900">{quotingTrip.pickup_zone}</span>
                                </div>
                                <div>
                                    <span className="text-gray-500 block text-xs uppercase tracking-wide">Dropoff</span>
                                    <span className="font-medium text-gray-900">{quotingTrip.dropoff_zone}</span>
                                </div>
                                <div>
                                    <span className="text-gray-500 block text-xs uppercase tracking-wide">Date & Time</span>
                                    <span className="font-medium text-gray-900">{new Date(quotingTrip.pickup_time).toLocaleString()}</span>
                                </div>
                                <div>
                                    <span className="text-gray-500 block text-xs uppercase tracking-wide">Vehicle</span>
                                    <span className="font-medium text-gray-900">{quotingTrip.requested_vehicle_type}</span>
                                </div>
                            </div>
                        </div>

                        <form onSubmit={handleSubmitQuote} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Price (AED)</label>
                                <input
                                    type="number"
                                    required
                                    min="0"
                                    step="0.01"
                                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                                    value={quotePrice}
                                    onChange={(e) => setQuotePrice(e.target.value)}
                                />
                            </div>
                            <div>
                                <label className="block text-sm font-medium text-gray-700 mb-1">Notes (Optional)</label>
                                <textarea
                                    className="w-full border border-gray-300 rounded-md px-3 py-2"
                                    rows={3}
                                    value={quoteNotes}
                                    onChange={(e) => setQuoteNotes(e.target.value)}
                                />
                            </div>
                            <div className="flex justify-end gap-2 pt-2">
                                <button
                                    type="button"
                                    onClick={() => setQuotingTrip(null)}
                                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                                >
                                    Cancel
                                </button>
                                <button
                                    type="submit"
                                    className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700"
                                >
                                    Submit
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default Dashboard;
