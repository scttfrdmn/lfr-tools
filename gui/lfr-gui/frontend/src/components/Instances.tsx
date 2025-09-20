import React, { useState, useEffect } from 'react'
import {
  Container,
  Header,
  Table,
  Box,
  SpaceBetween,
  Button,
  StatusIndicator,
  Badge,
  TextFilter,
  Pagination,
  CollectionPreferences,
  Modal,
  Alert
} from '@cloudscape-design/components'

import { LFRService } from "../../bindings/lfr-gui/pkg/services"

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

interface InstanceInfo {
  name: string
  state: string
  public_ip: string
  blueprint: string
  bundle: string
  region: string
  tags: Record<string, string>
  username: string
}

interface InstancesProps {
  userInfo: UserInfo
}

const Instances: React.FC<InstancesProps> = ({ userInfo }) => {
  const [instances, setInstances] = useState<InstanceInfo[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedItems, setSelectedItems] = useState<InstanceInfo[]>([])
  const [filteringText, setFilteringText] = useState('')
  const [currentPageIndex, setCurrentPageIndex] = useState(1)
  const [showStartModal, setShowStartModal] = useState(false)
  const [showStopModal, setShowStopModal] = useState(false)
  const [operationLoading, setOperationLoading] = useState(false)

  useEffect(() => {
    loadInstances()
  }, [userInfo.project])

  const loadInstances = async () => {
    setLoading(true)
    try {
      const instanceList = await LFRService.ListInstances(userInfo.project)
      setInstances(instanceList)
    } catch (err) {
      console.error('Failed to load instances:', err)
    }
    setLoading(false)
  }

  const handleStartInstances = async () => {
    setOperationLoading(true)
    try {
      for (const instance of selectedItems) {
        await LFRService.StartInstance(instance.name)
      }
      setShowStartModal(false)
      setSelectedItems([])
      // Refresh the list
      setTimeout(loadInstances, 2000) // Give instances time to start
    } catch (err) {
      console.error('Failed to start instances:', err)
    }
    setOperationLoading(false)
  }

  const handleStopInstances = async () => {
    setOperationLoading(true)
    try {
      for (const instance of selectedItems) {
        await LFRService.StopInstance(instance.name)
      }
      setShowStopModal(false)
      setSelectedItems([])
      // Refresh the list
      setTimeout(loadInstances, 2000) // Give instances time to stop
    } catch (err) {
      console.error('Failed to stop instances:', err)
    }
    setOperationLoading(false)
  }

  const getStateIndicator = (state: string) => {
    switch (state.toLowerCase()) {
      case 'running':
        return <StatusIndicator type="success">Running</StatusIndicator>
      case 'stopped':
        return <StatusIndicator type="stopped">Stopped</StatusIndicator>
      case 'pending':
      case 'starting':
        return <StatusIndicator type="in-progress">Starting</StatusIndicator>
      case 'stopping':
        return <StatusIndicator type="in-progress">Stopping</StatusIndicator>
      default:
        return <StatusIndicator type="info">{state}</StatusIndicator>
    }
  }

  const getBundleBadge = (bundle: string) => {
    if (bundle.includes('gpu')) {
      return <Badge color="red">GPU</Badge>
    } else if (bundle.includes('4xl')) {
      return <Badge color="orange">4XL</Badge>
    } else if (bundle.includes('2xl')) {
      return <Badge color="blue">2XL</Badge>
    } else {
      return <Badge color="green">XL</Badge>
    }
  }

  const filteredInstances = instances.filter(instance =>
    !filteringText ||
    instance.name.toLowerCase().includes(filteringText.toLowerCase()) ||
    instance.username.toLowerCase().includes(filteringText.toLowerCase()) ||
    instance.state.toLowerCase().includes(filteringText.toLowerCase())
  )

  const pageSize = 10
  const paginatedInstances = filteredInstances.slice(
    (currentPageIndex - 1) * pageSize,
    currentPageIndex * pageSize
  )

  return (
    <Container
      header={
        <Header
          variant="h1"
          actions={
            <SpaceBetween direction="horizontal" size="xs">
              <Button
                iconName="refresh"
                onClick={loadInstances}
                loading={loading}
              >
                Refresh
              </Button>
              {selectedItems.length > 0 && (
                <>
                  <Button
                    iconName="upload"
                    onClick={() => setShowStartModal(true)}
                    disabled={selectedItems.every(i => i.state === 'running')}
                  >
                    Start Selected
                  </Button>
                  <Button
                    iconName="download"
                    onClick={() => setShowStopModal(true)}
                    disabled={selectedItems.every(i => i.state === 'stopped')}
                  >
                    Stop Selected
                  </Button>
                </>
              )}
            </SpaceBetween>
          }
        >
          {userInfo.role === 'student' ? 'My Cloud Computer' : 'Class Instances'}
        </Header>
      }
    >
      <Table
        columnDefinitions={[
          {
            id: "name",
            header: "Instance",
            cell: (item: InstanceInfo) => item.name,
            sortingField: "name",
            isRowHeader: true
          },
          {
            id: "username",
            header: "Student",
            cell: (item: InstanceInfo) => item.username,
            sortingField: "username"
          },
          {
            id: "state",
            header: "Status",
            cell: (item: InstanceInfo) => getStateIndicator(item.state),
            sortingField: "state"
          },
          {
            id: "public_ip",
            header: "Public IP",
            cell: (item: InstanceInfo) => item.public_ip || "-",
            sortingField: "public_ip"
          },
          {
            id: "bundle",
            header: "Size",
            cell: (item: InstanceInfo) => (
              <SpaceBetween direction="horizontal" size="xs">
                {getBundleBadge(item.bundle)}
                <span>{item.bundle.replace('app_standard_', '').replace('_1_0', '').toUpperCase()}</span>
              </SpaceBetween>
            ),
            sortingField: "bundle"
          },
          {
            id: "actions",
            header: "Actions",
            cell: (item: InstanceInfo) => (
              <SpaceBetween direction="horizontal" size="xs">
                <Button
                  variant="link"
                  iconName="external"
                  disabled={item.state !== 'running'}
                >
                  SSH
                </Button>
                {item.state === 'running' ? (
                  <Button
                    variant="link"
                    iconName="download"
                    onClick={() => {
                      setSelectedItems([item])
                      setShowStopModal(true)
                    }}
                  >
                    Stop
                  </Button>
                ) : (
                  <Button
                    variant="link"
                    iconName="upload"
                    onClick={() => {
                      setSelectedItems([item])
                      setShowStartModal(true)
                    }}
                  >
                    Start
                  </Button>
                )}
              </SpaceBetween>
            )
          }
        ]}
        items={paginatedInstances}
        loading={loading}
        loadingText="Loading instances..."
        selectionType={userInfo.role !== 'student' ? "multi" : undefined}
        selectedItems={selectedItems}
        onSelectionChange={({ detail }) => setSelectedItems(detail.selectedItems)}
        header={
          <Header
            counter={`(${filteredInstances.length})`}
            actions={
              userInfo.role === 'student' ? undefined : (
                <SpaceBetween direction="horizontal" size="xs">
                  <Button iconName="add-plus">Add Instance</Button>
                  <Button iconName="settings">Bulk Actions</Button>
                </SpaceBetween>
              )
            }
          >
            Instances
          </Header>
        }
        filter={
          <TextFilter
            filteringText={filteringText}
            onChange={({ detail }) => setFilteringText(detail.filteringText)}
            filteringPlaceholder="Search instances..."
          />
        }
        pagination={
          <Pagination
            currentPageIndex={currentPageIndex}
            onChange={({ detail }) => setCurrentPageIndex(detail.currentPageIndex)}
            pagesCount={Math.ceil(filteredInstances.length / pageSize)}
          />
        }
        preferences={
          <CollectionPreferences
            title="Preferences"
            confirmLabel="Confirm"
            cancelLabel="Cancel"
            preferences={{
              pageSize: pageSize,
              visibleContent: ['name', 'username', 'state', 'public_ip', 'bundle', 'actions']
            }}
          />
        }
      />

      {/* Start Confirmation Modal */}
      <Modal
        onDismiss={() => setShowStartModal(false)}
        visible={showStartModal}
        header="Start Instances"
        footer={
          <Box float="right">
            <SpaceBetween direction="horizontal" size="xs">
              <Button variant="link" onClick={() => setShowStartModal(false)}>
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handleStartInstances}
                loading={operationLoading}
              >
                Start {selectedItems.length} Instance{selectedItems.length !== 1 ? 's' : ''}
              </Button>
            </SpaceBetween>
          </Box>
        }
      >
        <SpaceBetween direction="vertical" size="m">
          <Alert statusIconAriaLabel="Info" type="info">
            Starting instances will begin billing charges. Instances will be available for SSH connection within 1-2 minutes.
          </Alert>
          <Box>
            <strong>Instances to start:</strong>
            <ul>
              {selectedItems.map(instance => (
                <li key={instance.name}>{instance.name} ({instance.username})</li>
              ))}
            </ul>
          </Box>
        </SpaceBetween>
      </Modal>

      {/* Stop Confirmation Modal */}
      <Modal
        onDismiss={() => setShowStopModal(false)}
        visible={showStopModal}
        header="Stop Instances"
        footer={
          <Box float="right">
            <SpaceBetween direction="horizontal" size="xs">
              <Button variant="link" onClick={() => setShowStopModal(false)}>
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handleStopInstances}
                loading={operationLoading}
              >
                Stop {selectedItems.length} Instance{selectedItems.length !== 1 ? 's' : ''}
              </Button>
            </SpaceBetween>
          </Box>
        }
      >
        <SpaceBetween direction="vertical" size="m">
          <Alert statusIconAriaLabel="Warning" type="warning">
            Stopping instances will disconnect users and may cause data loss if work is not saved.
            Stopped instances will not incur charges.
          </Alert>
          <Box>
            <strong>Instances to stop:</strong>
            <ul>
              {selectedItems.map(instance => (
                <li key={instance.name}>{instance.name} ({instance.username})</li>
              ))}
            </ul>
          </Box>
        </SpaceBetween>
      </Modal>
    </Container>
  )
}

export default Instances