import React from 'react';
import { NavLink, Outlet, useNavigate } from 'react-router-dom';
import {
    LayoutDashboard,
    FileText,
    Bell,
    Calendar,
    LogOut,
    Menu,
    X,
    Building,
    Car,
    User,
    Users,
    DollarSign,
    Shield,
    Server
} from 'lucide-react';
import { useAuthStore } from '../lib/auth';
import ChangePasswordModal from './ChangePasswordModal';

const Layout = () => {
    const navigate = useNavigate();
    const { logout, user, isAuthenticated } = useAuthStore();
    const [isMobileMenuOpen, setIsMobileMenuOpen] = React.useState(false);

    React.useEffect(() => {
        if (!isAuthenticated) {
            console.warn('Layout: Not authenticated, redirecting to login');
            navigate('/login');
        }
    }, [isAuthenticated, navigate]);

    const handleLogout = () => {
        logout();
        navigate('/login');
    };

    // Define Nav Items based on Role
    const getNavItems = () => {
        const role = user?.role || '';

        // Super Admin
        if (role === 'SUPER_ADMIN') {
            return [
                { name: 'System Dashboard', path: '/', icon: Shield },
                { name: 'Infrastructure', path: '/superadmin/dashboard', icon: Server },
            ];
        }

        // Operator / Admin
        if (role === 'ADMIN' || role === 'OPS' || role === 'DISPATCHER') {
            const items = [
                { name: 'Dashboard', path: '/', icon: LayoutDashboard },
                { name: 'Trips', path: '/dispatch', icon: LayoutDashboard }, // Path remains same for now to avoid refactoring everything
                { name: 'Outsource Companies', path: '/outsource-master', icon: Building }, // New route
                { name: 'Submitted Quotes', path: '/quotes', icon: FileText },
            ];

            if (role === 'ADMIN') {
                items.splice(2, 0, { name: 'Team Management', path: '/team', icon: Users });
            }

            return items;
        }

        // Partner / Supplier
        // Default to Partner view
        return [
            { name: 'Dashboard', path: '/partner/dashboard', icon: LayoutDashboard },
            { name: 'RFQs', path: '/partner/rfqs', icon: Bell },
            { name: 'Assigned Trips', path: '/partner/trips', icon: Calendar },
            { name: 'Vehicles', path: '/partner/vehicles', icon: Car },
            { name: 'Drivers', path: '/partner/drivers', icon: User },
            { name: 'Invoices', path: '/partner/invoices', icon: DollarSign },
        ];
    };

    const navItems = getNavItems();

    return (
        <div className="flex h-screen bg-gray-50">
            {/* Force Password Change Modal */}
            {user && (
                <ChangePasswordModal
                    isOpen={user.password_change_required || false}
                    userId={user.id} // User ID needed in AuthStore!
                    force={true}
                />
            )}

            {/* Sidebar for Desktop */}
            <aside className="hidden md:flex md:w-64 md:flex-col bg-white border-r border-gray-200">
                <div className="flex items-center justify-center h-16 border-b border-gray-200 px-6">
                    <h1 className="text-xl font-bold bg-gradient-to-r from-emerald-500 to-teal-600 text-transparent bg-clip-text">
                        Limo XL Connect
                    </h1>
                </div>
                <nav className="flex-1 px-4 py-4 space-y-1">
                    {navItems.map((item) => (
                        <NavLink
                            key={item.path}
                            to={item.path}
                            end={item.path === '/'}
                            className={({ isActive }) =>
                                `flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-lg transition-colors ${isActive
                                    ? 'bg-emerald-50 text-emerald-700'
                                    : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                                }`
                            }
                        >
                            <item.icon className="h-5 w-5" />
                            {item.name}
                        </NavLink>
                    ))}
                </nav>
                <div className="p-4 border-t border-gray-200">
                    <div className="mb-4 px-2">
                        <p className="text-xs text-gray-500 uppercase font-semibold">Logged in as</p>
                        <p className="text-sm font-bold text-gray-800 truncate">{user?.name}</p>
                        <p className="text-xs text-gray-500 truncate">{user?.email}</p>
                    </div>
                    <button
                        onClick={handleLogout}
                        className="flex items-center gap-3 w-full px-4 py-2 text-sm font-medium text-gray-600 hover:text-red-600 transition-colors"
                    >
                        <LogOut className="h-5 w-5" />
                        Sign Out
                    </button>
                </div>
            </aside>

            {/* Mobile Header & Content */}
            <div className="flex-1 flex flex-col overflow-hidden">
                <header className="md:hidden flex items-center justify-between px-4 py-3 bg-white border-b border-gray-200">
                    <span className="text-lg font-bold text-gray-800">Limo XL Connect</span>
                    <button onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}>
                        {isMobileMenuOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
                    </button>
                </header>

                {/* Mobile Menu */}
                {isMobileMenuOpen && (
                    <div className="md:hidden absolute top-14 left-0 w-full bg-white border-b border-gray-200 shadow-lg z-50">
                        <nav className="p-4 space-y-2">
                            {navItems.map((item) => (
                                <NavLink
                                    key={item.path}
                                    to={item.path}
                                    onClick={() => setIsMobileMenuOpen(false)}
                                    className={({ isActive }) =>
                                        `flex items-center gap-3 px-4 py-3 rounded-md ${isActive ? 'bg-emerald-50 text-emerald-700' : 'text-gray-600'
                                        }`
                                    }
                                >
                                    <item.icon className="h-5 w-5" />
                                    {item.name}
                                </NavLink>
                            ))}
                            <button
                                onClick={handleLogout}
                                className="flex items-center gap-3 w-full px-4 py-3 text-red-600"
                            >
                                <LogOut className="h-5 w-5" />
                                Sign Out
                            </button>
                        </nav>
                    </div>
                )}

                <main className="flex-1 overflow-auto p-4 md:p-8">
                    <Outlet />
                </main>
            </div>
        </div>
    );
};

export default Layout;
