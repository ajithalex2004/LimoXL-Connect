import React, { useEffect, useState } from 'react';
import { TripOffer, operatorService, tripService, Trip } from '../../services/api';
import { DollarSign, Clock, X, MapPin, User, Tag, CheckCircle, Car } from 'lucide-react';

const SubmittedQuotes: React.FC = () => {
    const [quotes, setQuotes] = useState<TripOffer[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [selectedTrip, setSelectedTrip] = useState<Trip | null>(null);
    const [selectedQuoteId, setSelectedQuoteId] = useState<string | null>(null);
    const [isModalOpen, setIsModalOpen] = useState(false);
    const [isRejectModalOpen, setIsRejectModalOpen] = useState(false);
    const [rejectionReason, setRejectionReason] = useState('');

    useEffect(() => {
        fetchQuotes();
    }, []);

    const fetchQuotes = async () => {
        try {
            const data = await operatorService.getSubmittedQuotes();
            setQuotes(Array.isArray(data) ? data : []);
        } catch (err) {
            setError('Failed to load quotes');
            console.error(err);
        } finally {
            setLoading(false);
        }
    };

    const handleViewTrip = async (quote: TripOffer) => {
        try {
            setSelectedQuoteId(quote.id);
            const trip = await tripService.getTrip(quote.trip_id);
            setSelectedTrip(trip);
            setIsModalOpen(true);
        } catch (err) {
            console.error("Failed to load trip details", err);
            alert("Failed to load trip details");
        }
    };

    const handleAcceptQuote = async () => {
        if (!selectedQuoteId) return;
        if (!confirm('Are you sure you want to accept this quote? This will assign the trip to the supplier.')) return;

        try {
            await operatorService.acceptQuote(selectedQuoteId);
            alert('Quote Accepted!');
            closeModal();
            fetchQuotes(); // Refresh list to show updated status
        } catch (err) {
            console.error("Failed to accept quote", err);
            alert("Failed to accept quote");
        }
    };

    const closeModal = () => {
        setIsModalOpen(false);
        setSelectedTrip(null);
        setSelectedQuoteId(null);
    };

    const openRejectModal = () => {
        setIsRejectModalOpen(true);
        setRejectionReason('');
    };

    const closeRejectModal = () => {
        setIsRejectModalOpen(false);
        setRejectionReason('');
    };

    const handleRejectQuote = async () => {
        if (!selectedQuoteId) return;
        if (!rejectionReason.trim()) {
            alert('Please provide a reason for rejection');
            return;
        }

        try {
            await operatorService.rejectQuote(selectedQuoteId);
            alert(`Quote Rejected. Reason: ${rejectionReason}`);
            closeRejectModal();
            closeModal();
            fetchQuotes(); // Refresh list
        } catch (err) {
            console.error("Failed to reject quote", err);
            alert("Failed to reject quote");
        }
    };

    if (loading) return <div className="p-8 text-center">Loading quotes...</div>;
    if (error) return <div className="p-8 text-center text-red-500">{error}</div>;

    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold text-gray-900">Submitted Quotes</h1>

            <div className="bg-white shadow-sm rounded-lg overflow-hidden border border-gray-200">
                <table className="min-w-full divide-y divide-gray-200">
                    <thead className="bg-gray-50">
                        <tr>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created At</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">RFQ #</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Supplier</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Price</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Notes</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                            <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
                        </tr>
                    </thead>
                    <tbody className="bg-white divide-y divide-gray-200">
                        {(!quotes || quotes.length === 0) ? (
                            <tr>
                                <td colSpan={6} className="px-6 py-4 text-center text-gray-500">
                                    No quotes submitted yet.
                                </td>
                            </tr>
                        ) : (
                            quotes.map((quote) => (
                                <tr key={quote.id}>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                        <div className="flex items-center">
                                            {new Date(quote.created_at).toLocaleString()}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 font-mono">
                                        {quote.rfq_number || quote.booking_reference || 'N/A'}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                        {quote.supplier_name}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                        <div className="flex items-center text-green-600 font-semibold">
                                            <DollarSign className="h-4 w-4 mr-1" />
                                            {quote.price.toFixed(2)}
                                        </div>
                                    </td>
                                    <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">
                                        {quote.notes ? (
                                            <span title={quote.notes}>{quote.notes}</span>
                                        ) : (
                                            <span className="text-gray-400 italic">No notes</span>
                                        )}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-xs font-bold text-gray-500 uppercase">
                                        {quote.service_type || 'One Way'}
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap">
                                        <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${quote.status === 'ACCEPTED' ? 'bg-green-100 text-green-800' :
                                            quote.status === 'REJECTED' ? 'bg-red-100 text-red-800' :
                                                quote.status === 'TIMEOUT' ? 'bg-gray-100 text-gray-800' :
                                                    'bg-yellow-100 text-yellow-800'
                                            }`}>
                                            {quote.status}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 whitespace-nowrap text-sm text-blue-600 hover:text-blue-800 cursor-pointer font-medium" onClick={() => handleViewTrip(quote)}>
                                        View Details
                                    </td>
                                </tr>
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {/* Trip Details Modal */}
            {isModalOpen && selectedTrip && (
                <div className="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
                    <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
                        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" onClick={closeModal}></div>
                        <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
                        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
                            <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                                <div className="sm:flex sm:items-start">
                                    <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                                        <div className="flex justify-between items-center mb-4">
                                            <h3 className="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                                                Trip Details: {selectedTrip.booking_reference}
                                            </h3>
                                            <button onClick={closeModal} className="text-gray-400 hover:text-gray-500">
                                                <X className="h-6 w-6" />
                                            </button>
                                        </div>

                                        <div className="bg-gray-50 rounded-lg p-4 space-y-3">
                                            {/* Quoted Price */}
                                            {selectedQuoteId && quotes.find(q => q.id === selectedQuoteId) && (
                                                <div className="flex items-start gap-3 bg-white p-3 rounded-md border border-gray-200">
                                                    <DollarSign className="h-5 w-5 text-green-600 mt-0.5" />
                                                    <div>
                                                        <p className="text-sm font-medium text-gray-900">Your Quote</p>
                                                        <p className="text-lg font-bold text-green-600">
                                                            AED {quotes.find(q => q.id === selectedQuoteId)?.price.toFixed(2)}
                                                        </p>
                                                    </div>
                                                </div>
                                            )}

                                            <div className="flex items-start gap-3">
                                                <User className="h-5 w-5 text-gray-400 mt-0.5" />
                                                <div>
                                                    <p className="text-sm font-medium text-gray-900">Passenger</p>
                                                    <p className="text-sm text-gray-600">{selectedTrip.passenger_name}</p>
                                                </div>
                                            </div>

                                            <div className="flex items-start gap-3">
                                                <MapPin className="h-5 w-5 text-gray-400 mt-0.5" />
                                                <div className="flex-1">
                                                    <p className="text-sm font-medium text-gray-900">Route</p>
                                                    <div className="flex flex-col gap-1 mt-1">
                                                        <div className="flex items-center gap-2">
                                                            <div className="w-2 h-2 rounded-full bg-green-500"></div>
                                                            <span className="text-sm text-gray-600">{selectedTrip.pickup_zone}</span>
                                                        </div>
                                                        <div className="w-0.5 h-3 bg-gray-300 ml-0.5"></div>
                                                        <div className="flex items-center gap-2">
                                                            <div className="w-2 h-2 rounded-full bg-red-500"></div>
                                                            <span className="text-sm text-gray-600">{selectedTrip.dropoff_zone}</span>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>

                                            <div className="flex items-start gap-3">
                                                <Clock className="h-5 w-5 text-gray-400 mt-0.5" />
                                                <div>
                                                    <p className="text-sm font-medium text-gray-900">Pickup Time</p>
                                                    <p className="text-sm text-gray-600">{new Date(selectedTrip.pickup_time).toLocaleString()}</p>
                                                </div>
                                            </div>

                                            <div className="flex items-start gap-3">
                                                <Tag className="h-5 w-5 text-gray-400 mt-0.5" />
                                                <div>
                                                    <p className="text-sm font-medium text-gray-900">Request Type</p>
                                                    <p className="text-sm text-gray-600 font-bold uppercase">{selectedTrip.service_type || 'One Way'}</p>
                                                </div>
                                            </div>

                                            <div className="flex items-start gap-3">
                                                <Car className="h-5 w-5 text-gray-400 mt-0.5" />
                                                <div>
                                                    <p className="text-sm font-medium text-gray-900">Vehicle Requested</p>
                                                    <p className="text-sm text-gray-600">
                                                        {selectedTrip.requested_vehicle_type}
                                                        {selectedTrip.requested_vehicle_class ? ` (${selectedTrip.requested_vehicle_class})` : ''}
                                                    </p>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-2">
                                {/* Only show Accept button if pending */}
                                {quotes.find(q => q.id === selectedQuoteId)?.status === 'PENDING' && (
                                    <>
                                        <button type="button" className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-emerald-600 text-base font-medium text-white hover:bg-emerald-700 focus:outline-none sm:w-auto sm:text-sm" onClick={handleAcceptQuote}>
                                            <CheckCircle className="h-5 w-5 mr-2" />
                                            Accept Quote
                                        </button>
                                        <button type="button" className="mt-3 w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" onClick={openRejectModal}>
                                            <X className="h-5 w-5 mr-2" />
                                            Reject Quote
                                        </button>
                                    </>
                                )}

                                <button type="button" className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm" onClick={closeModal}>
                                    Close
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}

            {/* Rejection Reason Modal */}
            {isRejectModalOpen && (
                <div className="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="reject-modal-title" role="dialog" aria-modal="true">
                    <div className="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
                        <div className="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" onClick={closeRejectModal}></div>
                        <span className="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
                        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
                            <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
                                <div className="sm:flex sm:items-start">
                                    <div className="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                                        <X className="h-6 w-6 text-red-600" />
                                    </div>
                                    <div className="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left w-full">
                                        <h3 className="text-lg leading-6 font-medium text-gray-900" id="reject-modal-title">
                                            Reject Quote
                                        </h3>
                                        <div className="mt-4">
                                            <label htmlFor="rejection-reason" className="block text-sm font-medium text-gray-700 mb-2">
                                                Reason for Rejection *
                                            </label>
                                            <textarea
                                                id="rejection-reason"
                                                rows={4}
                                                className="shadow-sm focus:ring-red-500 focus:border-red-500 block w-full sm:text-sm border-gray-300 rounded-md"
                                                placeholder="Please provide a reason for rejecting this quote..."
                                                value={rejectionReason}
                                                onChange={(e) => setRejectionReason(e.target.value)}
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                            <div className="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
                                <button
                                    type="button"
                                    className="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none sm:ml-3 sm:w-auto sm:text-sm"
                                    onClick={handleRejectQuote}
                                >
                                    Confirm Rejection
                                </button>
                                <button
                                    type="button"
                                    className="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none sm:mt-0 sm:w-auto sm:text-sm"
                                    onClick={closeRejectModal}
                                >
                                    Cancel
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default SubmittedQuotes;
