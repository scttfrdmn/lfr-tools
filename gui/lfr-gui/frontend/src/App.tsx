import { useState, useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import '@cloudscape-design/global-styles/index.css'
import { AppLayout, TopNavigation, SideNavigation } from '@cloudscape-design/components'

import { LFRService } from "../bindings/lfr-gui/pkg/services"

// Import components
import Dashboard from './components/Dashboard'
import Instances from './components/Instances'
import Users from './components/Users'
import Students from './components/Students'
import Settings from './components/Settings'
import Analytics from './components/Analytics'

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

function App() {
  const [userInfo, setUserInfo] = useState<UserInfo | null>(null)
  const [loading, setLoading] = useState(true)
  const [activeHref, setActiveHref] = useState('/dashboard')

  useEffect(() => {
    // Get user role and permissions on app load
    LFRService.GetUserRole()
      .then((user: UserInfo | null) => {
        setUserInfo(user)
        setLoading(false)
      })
      .catch((err: any) => {
        console.error('Failed to get user role:', err)
        setLoading(false)
      })
  }, [])

  // Navigation items based on user role
  const getNavigationItems = () => {
    if (!userInfo) return []

    const baseItems = [
      { type: 'link', text: 'Dashboard', href: '/dashboard' },
      { type: 'link', text: 'Instances', href: '/instances' },
    ]

    // Add role-specific navigation
    if (userInfo.role === 'professor' || userInfo.role === 'admin') {
      baseItems.push(
        { type: 'link', text: 'Users & Groups', href: '/users' },
        { type: 'link', text: 'Students', href: '/students' },
        { type: 'link', text: 'Analytics', href: '/analytics' },
        { type: 'link', text: 'Settings', href: '/settings' }
      )
    } else if (userInfo.role === 'ta') {
      baseItems.push(
        { type: 'link', text: 'Student Support', href: '/students' },
        { type: 'link', text: 'Class Status', href: '/analytics' }
      )
    }

    return baseItems
  }

  const topNavigation = (
    <TopNavigation
      identity={{
        href: "/",
        title: "LFR Tools"
      }}
      utilities={[
        {
          type: "menu-dropdown",
          text: userInfo?.username || "User",
          description: userInfo?.role || "Loading...",
          iconName: "user-profile",
          items: [
            { id: "profile", text: "Profile" },
            { id: "settings", text: "Settings" },
            { id: "signout", text: "Sign out" }
          ]
        }
      ]}
    />
  )

  const sideNavigation = (
    <SideNavigation
      activeHref={activeHref}
      header={{ href: "/dashboard", text: userInfo?.project || "Project" }}
      items={getNavigationItems()}
      onFollow={(event) => {
        if (!event.detail.external) {
          setActiveHref(event.detail.href)
        }
      }}
    />
  )

  if (loading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        backgroundColor: '#f9f9f9'
      }}>
        <div>Loading LFR Tools...</div>
      </div>
    )
  }

  if (!userInfo) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        height: '100vh',
        backgroundColor: '#f9f9f9'
      }}>
        <div>Failed to load user information. Please check your configuration.</div>
      </div>
    )
  }

  return (
    <Router>
      {topNavigation}
      <AppLayout
        navigation={sideNavigation}
        content={
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<Dashboard userInfo={userInfo} />} />
            <Route path="/instances" element={<Instances userInfo={userInfo} />} />
            {(userInfo.role === 'professor' || userInfo.role === 'admin') && (
              <>
                <Route path="/users" element={<Users userInfo={userInfo} />} />
                <Route path="/students" element={<Students userInfo={userInfo} />} />
                <Route path="/analytics" element={<Analytics userInfo={userInfo} />} />
                <Route path="/settings" element={<Settings userInfo={userInfo} />} />
              </>
            )}
            {userInfo.role === 'ta' && (
              <>
                <Route path="/students" element={<Students userInfo={userInfo} />} />
                <Route path="/analytics" element={<Analytics userInfo={userInfo} />} />
              </>
            )}
          </Routes>
        }
        notifications={<div />} // Placeholder for notifications
        breadcrumbs={<div />}  // Placeholder for breadcrumbs
        disableContentHeaderOverlap={true}
      />
    </Router>
  )
}

export default App