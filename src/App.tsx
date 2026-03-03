import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import RFQs from './pages/RFQs';
import Trips from './pages/Trips';
import Invoices from './pages/Invoices';
import Vehicles from './pages/partner/Vehicles';
import Drivers from './pages/partner/Drivers';
import Login from './pages/Login';
import OutsourceMaster from './pages/operator/OutsourceMaster';
import SubmittedQuotes from './pages/operator/SubmittedQuotes';
import TeamManagement from './pages/operator/TeamManagement';
import Dispatch from './pages/operator/Dispatch';
import OperatorVehicles from './pages/operator/Vehicles';
import OperatorDrivers from './pages/operator/Drivers';
import SecureDriverLink from './pages/SecureDriverLink';
import TenantDashboard from './pages/superadmin/TenantDashboard';
import TenantDetail from './pages/superadmin/TenantDetail';

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Layout />}>
                    <Route index element={<Dashboard />} />
                    <Route path="partner/dashboard" element={<Dashboard />} />
                    <Route path="rfqs" element={<RFQs />} />
                    <Route path="partner/rfqs" element={<RFQs />} />
                    <Route path="partner/trips" element={<Trips />} />
                    <Route path="partner/vehicles" element={<Vehicles />} />
                    <Route path="partner/drivers" element={<Drivers />} />
                    <Route path="partner/invoices" element={<Invoices />} />
                    <Route path="outsource-master" element={<OutsourceMaster />} />
                    <Route path="dispatch" element={<Dispatch />} />
                    <Route path="quotes" element={<SubmittedQuotes />} />
                    <Route path="team" element={<TeamManagement />} />
                    <Route path="vehicles" element={<OperatorVehicles />} />
                    <Route path="drivers" element={<OperatorDrivers />} />
                    
                    {/* SuperAdmin Routes */}
                    <Route path="superadmin/dashboard" element={<TenantDashboard />} />
                    <Route path="superadmin/tenants/:id" element={<TenantDetail />} />
                </Route>
                <Route path="/login" element={<Login />} />
                <Route path="/link/:token" element={<SecureDriverLink />} />
                <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;
