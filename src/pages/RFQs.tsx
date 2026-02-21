import { useState, useEffect } from 'react';
import { partnerService, type Trip } from '../services/api';
import {
    MapPin,
    Clock,
    CheckCircle,
    Car,
    Briefcase
} from 'lucide-react';

const RFQs = () => {
    const [activeTab, setActiveTab] = useState<'new' | 'submitted' | 'history'>('new');
    const [rfqs, setRfqs] = useState<Trip[]>([]);
    const [historyTrips, setHistoryTrips] = useState<Trip[]>([]);
    const [loading, setLoading] = useState(true);
    const [submitting, setSubmitting] = useState<string | null>(null);
    const [quoteInputs, setQuoteInputs] = useState<{ [key: string]: string }>({});

    useEffect(() => {
        loadRFQs();
    }, []);

    const loadRFQs = async () => {
        try {
            setLoading(true);
            const [openRfqs, history] = await Promise.all([
                partnerService.listRFQs(),
                partnerService.listRFQHistory()
            ]);
            setRfqs(openRfqs || []);
            setHistoryTrips(history || []);
        } catch (error) {
            console.error("Failed to load RFQs", error);
        } finally {
            setLoading(false);
        }
    };

    const handleQuoteChange = (tripId: string, value: string) => {
        setQuoteInputs(prev => ({ ...prev, [tripId]: value }));
    };

    const formatTripDate = (dateStr: string) => {
        try {
            const d = new Date(dateStr);
            if (isNaN(d.getTime())) return "Invalid Date";
            return `${d.toLocaleDateString()} at ${d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
        } catch (e) {
            return "Date Error";
        }
    };

    const handleQuoteSubmit = async (tripId: string) => {
        const amountStr = quoteInputs[tripId];
        if (!amountStr) {
            alert("Please enter a quote amount");
            return;
        }

        const price = parseFloat(amountStr);
        if (isNaN(price) || price <= 0) {
            alert("Please enter a valid amount");
            return;
        }

        setSubmitting(tripId);
        try {
            await partnerService.submitQuote({ trip_id: tripId, price });
            alert("Quote submitted successfully!");
            // In a real app we might remove the item from list or mark as quoted
            // For now, reload
            loadRFQs();
            setQuoteInputs(prev => ({ ...prev, [tripId]: '' }));
        } catch (error) {
            console.error("Failed to submit quote", error);
            alert("Failed to submit quote");
        } finally {
            setSubmitting(null);
        }
    };

    if (loading) return <div className="p-8 text-center text-gray-500">Loading RFQs...</div>;

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">RFQ Management</h1>
                    <p className="text-gray-500">Manage incoming requests and submitted quotes</p>
                </div>
            </div>

            <div className="border-b border-gray-200">
                <nav className="-mb-px flex space-x-8">
                    <button
                        onClick={() => setActiveTab('new')}
                        className={`${activeTab === 'new'
                            ? 'border-emerald-500 text-emerald-600'
                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
                    >
                        New Requests
                        <span className="bg-emerald-100 text-emerald-600 py-0.5 px-2 rounded-full text-xs">
                            {rfqs.filter(t => t.status === 'MARKETPLACE_SEARCH').length}
                        </span>
                    </button>
                    <button
                        onClick={() => setActiveTab('submitted')}
                        className={`${activeTab === 'submitted'
                            ? 'border-blue-500 text-blue-600'
                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
                    >
                        Submitted Quotes
                        <span className="bg-blue-100 text-blue-600 py-0.5 px-2 rounded-full text-xs">
                            {rfqs.filter(t => t.status === 'OFFERED').length}
                        </span>
                    </button>
                    <button
                        onClick={() => setActiveTab('history')}
                        className={`${activeTab === 'history'
                            ? 'border-gray-500 text-gray-600'
                            : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm flex items-center gap-2`}
                    >
                        History
                        <span className="bg-gray-100 text-gray-600 py-0.5 px-2 rounded-full text-xs">
                            {historyTrips.length + rfqs.filter(t => t.status === 'OFFERED').length}
                        </span>
                    </button>
                </nav>
            </div>

            <div className="grid gap-4">
                {(activeTab === 'new' ? rfqs.filter(t => t.status === 'MARKETPLACE_SEARCH') :
                    activeTab === 'submitted' ? rfqs.filter(t => t.status === 'OFFERED') :
                        // History: Combine rejected trips AND submitted quotes
                        [...historyTrips, ...rfqs.filter(t => t.status === 'OFFERED')]
                ).map((trip) => (
                    <div key={trip.id} className="bg-white rounded-lg shadow-sm border border-gray-200 p-6 flex flex-col lg:flex-row justify-between gap-6 hover:border-emerald-200 transition-colors relative overflow-hidden">
                        {/* Status Badge for History/Submitted */}
                        {activeTab !== 'new' && (
                            <div className={`absolute top-0 left-0 px-3 py-1 rounded-br-lg text-xs font-bold uppercase tracking-wide border-b border-r ${trip.status === 'REJECTED' ? 'bg-red-100 text-red-600 border-red-200' :
                                trip.status === 'OFFERED' ? 'bg-blue-100 text-blue-600 border-blue-200' :
                                    'bg-gray-100 text-gray-600 border-gray-200'
                                }`}>
                                {trip.status}
                            </div>
                        )}
                        {/* Service Type Badge */}
                        <div className="absolute top-0 right-0 bg-slate-100 text-slate-600 px-3 py-1 rounded-bl-lg text-xs font-bold uppercase tracking-wide border-b border-l border-slate-200">
                            {trip.service_type?.replace('_', ' ') || 'ONE WAY'}
                        </div>

                        <div className="flex-1 space-y-4 pt-2">
                            <div className="flex items-center gap-3">
                                <span className="bg-blue-100 text-blue-800 text-xs font-semibold px-2.5 py-0.5 rounded">
                                    {trip.booking_reference}
                                </span>
                                <span className="text-sm text-gray-500 flex items-center gap-1">
                                    <Clock className="h-3 w-3" />
                                    Posted recently
                                </span>
                            </div>

                            <div className="grid md:grid-cols-2 gap-6">
                                <div className="space-y-1">
                                    <div className="text-xs text-gray-400 uppercase tracking-wider font-semibold">Pickup</div>
                                    <div className="flex items-start gap-2 text-gray-900">
                                        <MapPin className="h-5 w-5 text-green-600 mt-0.5 shrink-0" />
                                        <div>
                                            <p className="font-medium text-lg">{trip.pickup_zone}</p>
                                            <p className="text-sm text-gray-500">
                                                {formatTripDate(trip.pickup_time)}
                                            </p>
                                        </div>
                                    </div>
                                </div>

                                <div className="space-y-1">
                                    <div className="text-xs text-gray-400 uppercase tracking-wider font-semibold">Dropoff</div>
                                    <div className="flex items-start gap-2 text-gray-900">
                                        <MapPin className="h-5 w-5 text-red-600 mt-0.5 shrink-0" />
                                        <div>
                                            <p className="font-medium text-lg">{trip.dropoff_zone}</p>
                                        </div>
                                    </div>
                                </div>
                            </div>

                            <div className="flex flex-wrap items-center gap-4 pt-2">
                                <div className="flex items-center gap-2 text-sm text-gray-700 bg-gray-50 px-3 py-1.5 rounded-lg border border-gray-200">
                                    <Car className="h-4 w-4 text-gray-500" />
                                    <span className="font-medium">
                                        {trip.requested_vehicle_class} {trip.requested_vehicle_group}
                                    </span>
                                </div>
                                <div className="flex items-center gap-2 text-sm text-gray-700 bg-gray-50 px-3 py-1.5 rounded-lg border border-gray-200">
                                    <Car className="h-4 w-4 text-gray-500" />
                                    <span className="font-medium">{trip.requested_vehicle_type || 'General'}</span>
                                </div>
                                <div className="flex items-center gap-2 text-sm text-gray-700 bg-gray-50 px-3 py-1.5 rounded-lg border border-gray-200">
                                    <Briefcase className="h-4 w-4 text-gray-500" />
                                    <span className="font-medium">{trip.service_type?.replace('_', ' ') || 'Standard'}</span>
                                </div>
                            </div>
                        </div>

                        {/* Quote Action Section */}
                        {/* Quote Action Section - Only for New Requests */}
                        {activeTab === 'new' && (
                            <div className="flex flex-col justify-end lg:w-72 border-t lg:border-t-0 lg:border-l border-gray-100 lg:pl-6 pt-4 lg:pt-0 gap-3">
                                <label className="text-sm font-medium text-gray-700">Enter Quote Amount (AED)</label>
                                <div className="flex gap-2">
                                    <div className="relative flex-1">
                                        <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                            <span className="text-gray-500 font-medium text-sm">AED</span>
                                        </div>
                                        <input
                                            type="number"
                                            placeholder="0.00"
                                            value={quoteInputs[trip.id] || ''}
                                            onChange={(e) => handleQuoteChange(trip.id, e.target.value)}
                                            className="block w-full pl-12 pr-3 py-2.5 border border-gray-300 rounded-lg focus:ring-emerald-500 focus:border-emerald-500 sm:text-sm"
                                        />
                                    </div>
                                    <button
                                        onClick={() => handleQuoteSubmit(trip.id)}
                                        disabled={submitting === trip.id}
                                        className="bg-emerald-600 hover:bg-emerald-700 text-white font-medium py-2 px-4 rounded-lg transition-colors disabled:opacity-50 whitespace-nowrap"
                                    >
                                        {submitting === trip.id ? 'Sending...' : 'Submit'}
                                    </button>
                                </div>
                                <p className="text-xs text-gray-400 text-center">
                                    Price includes all fees and taxes.
                                </p>
                            </div>
                        )}
                        {/* Info for Submitted/History */}
                        {activeTab !== 'new' && (
                            <div className="flex flex-col justify-center lg:w-auto border-t lg:border-t-0 lg:border-l border-gray-100 lg:pl-6 pt-4 lg:pt-0">
                                <div className="text-sm text-gray-500">
                                    {activeTab === 'submitted' ? (
                                        <div className="flex items-center gap-2 text-blue-600 font-medium">
                                            <Clock className="h-4 w-4" />
                                            Quote Submitted
                                        </div>
                                    ) : (
                                        <div className="flex items-center gap-2 text-gray-600 font-medium">
                                            {trip.status === 'REJECTED' ? (
                                                <span className="text-red-600 flex items-center gap-2">Quote Rejected</span>
                                            ) : (
                                                <span>{trip.status}</span>
                                            )}
                                        </div>
                                    )}
                                </div>
                            </div>
                        )}
                    </div>
                ))}

                {(
                    (activeTab === 'new' && rfqs.filter(t => t.status === 'MARKETPLACE_SEARCH').length === 0) ||
                    (activeTab === 'submitted' && rfqs.filter(t => t.status === 'OFFERED').length === 0) ||
                    (activeTab === 'history' && historyTrips.length + rfqs.filter(t => t.status === 'OFFERED').length === 0)
                ) && (
                        <div className="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-200">
                            <CheckCircle className="h-12 w-12 text-gray-300 mx-auto mb-4" />
                            <h3 className="text-lg font-medium text-gray-900">All caught up!</h3>
                            <p className="text-gray-500">No open RFQs available at the moment.</p>
                            <p className="text-xs text-gray-400 mt-2">Last checked: {new Date().toLocaleTimeString()}</p>
                        </div>
                    )}
            </div>
        </div>
    );
};

export default RFQs;
