import { useState, useEffect, useCallback } from 'react'
import { Events } from '@wailsio/runtime'

export interface StatusUpdate {
  type: string
  timestamp: string
  project: string
  data: any
}

export interface InstanceStatusUpdate {
  name: string
  state: string
  public_ip: string
  username: string
  previous_state: string
}

export interface ProjectStatusUpdate {
  name: string
  students_online: number
  budget_used: number
  alert_count: number
}

// Custom hook for real-time status monitoring
export const useRealTimeStatus = (project: string) => {
  const [lastUpdate, setLastUpdate] = useState<StatusUpdate | null>(null)
  const [projectStatus, setProjectStatus] = useState<ProjectStatusUpdate | null>(null)
  const [instanceChanges, setInstanceChanges] = useState<InstanceStatusUpdate[]>([])

  const handleStatusUpdate = useCallback((update: StatusUpdate) => {
    setLastUpdate(update)

    switch (update.type) {
      case 'project_status':
        setProjectStatus(update.data as ProjectStatusUpdate)
        break
      case 'instance_changes':
        setInstanceChanges(update.data as InstanceStatusUpdate[])
        break
      case 'instances_refresh':
        // Handle full instance list refresh
        break
    }
  }, [])

  useEffect(() => {
    if (!project) return

    // Subscribe to status updates for this project
    const eventName = `status_update_${project}`

    const unsubscribe = Events.On(eventName, (data: any) => {
      try {
        const update: StatusUpdate = JSON.parse(data.data)
        handleStatusUpdate(update)
      } catch (error) {
        console.error('Failed to parse status update:', error)
      }
    })

    // Request initial status
    // MonitoringService.SubscribeToProject(project)

    return () => {
      if (unsubscribe) {
        unsubscribe()
      }
    }
  }, [project, handleStatusUpdate])

  return {
    lastUpdate,
    projectStatus,
    instanceChanges,
    isConnected: !!lastUpdate
  }
}

// Hook for monitoring a specific student's instance
export const useStudentInstanceStatus = (username: string, project: string) => {
  const { instanceChanges } = useRealTimeStatus(project)
  const [studentInstance, setStudentInstance] = useState<InstanceStatusUpdate | null>(null)

  useEffect(() => {
    // Find updates for this specific student
    const studentUpdate = instanceChanges.find(change => change.username === username)
    if (studentUpdate) {
      setStudentInstance(studentUpdate)
    }
  }, [instanceChanges, username])

  return {
    instance: studentInstance,
    hasRecentUpdate: !!studentInstance
  }
}

// Hook for professor dashboard monitoring
export const useClassMonitoring = (project: string) => {
  const { projectStatus, instanceChanges, lastUpdate } = useRealTimeStatus(project)
  const [alerts, setAlerts] = useState<string[]>([])

  useEffect(() => {
    // Generate alerts based on status changes
    const newAlerts: string[] = []

    instanceChanges.forEach(change => {
      if (change.state === 'running' && change.previous_state === 'stopped') {
        newAlerts.push(`${change.username} connected to their computer`)
      } else if (change.state === 'stopped' && change.previous_state === 'running') {
        newAlerts.push(`${change.username} disconnected`)
      }
    })

    if (newAlerts.length > 0) {
      setAlerts(prev => [...newAlerts, ...prev].slice(0, 10)) // Keep last 10 alerts
    }
  }, [instanceChanges])

  return {
    projectStatus,
    instanceChanges,
    alerts,
    lastUpdate: lastUpdate?.timestamp,
    studentsOnline: projectStatus?.students_online || 0,
    budgetUsed: projectStatus?.budget_used || 0
  }
}