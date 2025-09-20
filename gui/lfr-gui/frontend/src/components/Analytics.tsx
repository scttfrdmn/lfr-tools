import React, { useState, useEffect } from 'react'
import {
  Container,
  Header,
  Grid,
  Box,
  SpaceBetween,
  LineChart,
  BarChart,
  PieChart,
  Cards,
  Badge,
  Button,
  Alert,
  ProgressBar,
  StatusIndicator
} from '@cloudscape-design/components'

import { useClassMonitoring } from '../hooks/useRealTimeStatus'
import { LFRService } from "../../bindings/lfr-gui/pkg/services"

interface UserInfo {
  role: string
  username: string
  project: string
  permissions: string[]
}

interface AnalyticsProps {
  userInfo: UserInfo
}

const Analytics: React.FC<AnalyticsProps> = ({ userInfo }) => {
  const { projectStatus, studentsOnline, budgetUsed, alerts } = useClassMonitoring(userInfo.project)
  const [costData, setCostData] = useState<any[]>([])
  const [usageData, setUsageData] = useState<any[]>([])
  const [optimizationSuggestions, setOptimizationSuggestions] = useState<any[]>([])

  useEffect(() => {
    // Simulate cost data
    setCostData([
      { x: 'Week 1', y: 45 },
      { x: 'Week 2', y: 52 },
      { x: 'Week 3', y: 48 },
      { x: 'Week 4', y: 38 },
      { x: 'Week 5', y: 42 },
      { x: 'Week 6', y: 35 },
    ])

    // Simulate usage data
    setUsageData([
      { name: 'alice', value: 85, description: 'High usage' },
      { name: 'bob', value: 45, description: 'Medium usage' },
      { name: 'charlie', value: 92, description: 'Very high usage' },
      { name: 'diana', value: 12, description: 'Low usage' },
      { name: 'emma', value: 33, description: 'Low usage' },
    ])

    // Simulate optimization suggestions
    setOptimizationSuggestions([
      {
        type: 'resize',
        instance: 'alice-ubuntu',
        current: '2XL',
        suggested: 'XL',
        savings: '$15/month',
        reason: 'Low CPU utilization (avg 8%)'
      },
      {
        type: 'idle',
        instance: 'bob-ubuntu',
        current: 'No idle detection',
        suggested: 'Educational conservative',
        savings: '$12/month',
        reason: 'Instance idle 60% of time'
      },
      {
        type: 'gpu',
        instance: 'diana-gpu',
        current: 'GPU enabled',
        suggested: 'Disable GPU',
        savings: '$25/month',
        reason: 'No GPU usage detected'
      }
    ])
  }, [])

  const renderCostChart = () => (
    <Container header={<Header variant="h2">💰 Weekly Cost Trend</Header>}>
      <LineChart
        series={[
          {
            title: "Weekly Spending",
            type: "line",
            data: costData
          }
        ]}
        xDomain={costData.map(d => d.x)}
        yDomain={[0, 60]}
        i18nStrings={{
          legendAriaLabel: "Legend",
          chartAriaRoleDescription: "Line chart showing weekly spending",
          xTickFormatter: (value) => value,
          yTickFormatter: (value) => `$${value}`
        }}
        ariaLabel="Weekly cost trend"
        height={200}
      />
    </Container>
  )

  const renderUsageChart = () => (
    <Container header={<Header variant="h2">📊 Student Usage Patterns</Header>}>
      <BarChart
        series={[
          {
            title: "Usage Hours",
            type: "bar",
            data: usageData
          }
        ]}
        xDomain={usageData.map(d => d.name)}
        yDomain={[0, 100]}
        i18nStrings={{
          legendAriaLabel: "Legend",
          chartAriaRoleDescription: "Bar chart showing student usage patterns",
          xTickFormatter: (value) => value,
          yTickFormatter: (value) => `${value}h`
        }}
        ariaLabel="Student usage patterns"
        height={200}
      />
    </Container>
  )

  const renderOptimizationSuggestions = () => (
    <Container
      header={
        <Header
          variant="h2"
          description="AI-powered recommendations to reduce costs while maintaining performance"
        >
          🚀 Cost Optimization Suggestions
        </Header>
      }
    >
      <Cards
        cardDefinition={{
          header: (item: any) => (
            <SpaceBetween direction="horizontal" size="xs">
              <Badge color={item.type === 'resize' ? 'blue' : item.type === 'gpu' ? 'red' : 'green'}>
                {item.type.toUpperCase()}
              </Badge>
              <Box fontWeight="bold">{item.instance}</Box>
            </SpaceBetween>
          ),
          sections: [
            {
              content: (item: any) => (
                <SpaceBetween direction="vertical" size="xs">
                  <Box>
                    <strong>Current:</strong> {item.current} → <strong>Suggested:</strong> {item.suggested}
                  </Box>
                  <Box color="text-status-success" fontWeight="bold">
                    💰 Save {item.savings}
                  </Box>
                  <Box color="text-body-secondary" fontSize="body-s">
                    {item.reason}
                  </Box>
                </SpaceBetween>
              )
            }
          ]
        }}
        items={optimizationSuggestions}
        loading={false}
        empty={
          <Box textAlign="center" color="inherit">
            <b>No optimization suggestions</b>
            <Box variant="p" color="inherit">
              Your instances are already well-optimized!
            </Box>
          </Box>
        }
        header={
          <Header
            actions={
              <Button variant="primary" iconName="external">
                Apply All Suggestions
              </Button>
            }
          >
            Recommendations ({optimizationSuggestions.length})
          </Header>
        }
      />
    </Container>
  )

  const renderBudgetOverview = () => (
    <Container header={<Header variant="h2">💳 Budget Overview</Header>}>
      <Grid gridDefinition={[{ colspan: 6 }, { colspan: 6 }]}>
        <SpaceBetween direction="vertical" size="m">
          <Box>
            <ProgressBar
              value={(budgetUsed / 500) * 100}
              description="Semester budget usage"
              additionalInfo={`$${budgetUsed} of $500 used`}
              label="Budget Progress"
            />
          </Box>

          <SpaceBetween direction="horizontal" size="m">
            <Box textAlign="center">
              <Box fontSize="heading-xl" fontWeight="bold" color="text-status-success">
                $159.50
              </Box>
              <Box color="text-body-secondary">Remaining</Box>
            </Box>
            <Box textAlign="center">
              <Box fontSize="heading-xl" fontWeight="bold" color="text-status-info">
                45
              </Box>
              <Box color="text-body-secondary">Days Left</Box>
            </Box>
          </SpaceBetween>
        </SpaceBetween>

        <SpaceBetween direction="vertical" size="s">
          <Box fontWeight="bold">💡 Budget Insights</Box>
          <Box color="text-body-secondary" fontSize="body-s">
            Current spending rate: $7.56/day
          </Box>
          <Box color="text-body-secondary" fontSize="body-s">
            Projected total: $480 (within budget)
          </Box>
          <Box color="text-status-success" fontSize="body-s">
            ✅ On track for semester completion
          </Box>
        </SpaceBetween>
      </Grid>
    </Container>
  )

  const renderRecentActivity = () => (
    <Container header={<Header variant="h2">📱 Recent Activity</Header>}>
      <SpaceBetween direction="vertical" size="s">
        {alerts.length > 0 ? (
          alerts.slice(0, 5).map((alert, index) => (
            <Box key={index} color="text-body-secondary" fontSize="body-s">
              • {alert}
            </Box>
          ))
        ) : (
          <Box color="text-body-secondary">No recent activity</Box>
        )}
      </SpaceBetween>
    </Container>
  )

  return (
    <Container header={<Header variant="h1">📊 Class Analytics - {userInfo.project}</Header>}>
      <SpaceBetween direction="vertical" size="l">
        {/* Real-time Status Bar */}
        <Alert
          statusIconAriaLabel="Info"
          type="info"
          header="Live Monitoring Active"
          action={<Button iconName="refresh">Refresh Now</Button>}
        >
          Real-time monitoring is active. Data updates every 30 seconds.
          Students online: {studentsOnline} • Last update: {new Date().toLocaleTimeString()}
        </Alert>

        {/* Key Metrics Grid */}
        <Grid gridDefinition={[{ colspan: 6 }, { colspan: 6 }]}>
          {renderBudgetOverview()}
          {renderRecentActivity()}
        </Grid>

        {/* Charts Grid */}
        <Grid gridDefinition={[{ colspan: 6 }, { colspan: 6 }]}>
          {renderCostChart()}
          {renderUsageChart()}
        </Grid>

        {/* Optimization Suggestions */}
        {renderOptimizationSuggestions()}

        {/* Live Student Status */}
        <Container header={<Header variant="h2">👥 Live Student Status</Header>}>
          <SpaceBetween direction="vertical" size="s">
            <Grid gridDefinition={[{ colspan: 3 }, { colspan: 3 }, { colspan: 3 }, { colspan: 3 }]}>
              <Box textAlign="center">
                <Box fontSize="heading-xl" fontWeight="bold" color="text-status-success">
                  {studentsOnline}
                </Box>
                <Box color="text-body-secondary">Online Now</Box>
              </Box>
              <Box textAlign="center">
                <Box fontSize="heading-xl" fontWeight="bold" color="text-status-info">
                  12
                </Box>
                <Box color="text-body-secondary">Working</Box>
              </Box>
              <Box textAlign="center">
                <Box fontSize="heading-xl" fontWeight="bold" color="text-status-warning">
                  3
                </Box>
                <Box color="text-body-secondary">Idle</Box>
              </Box>
              <Box textAlign="center">
                <Box fontSize="heading-xl" fontWeight="bold">
                  8
                </Box>
                <Box color="text-body-secondary">Offline</Box>
              </Box>
            </Grid>

            <Alert statusIconAriaLabel="Info" type="info">
              💡 Students working late? Consider extending lab hours or adjusting idle detection policies.
            </Alert>
          </SpaceBetween>
        </Container>
      </SpaceBetween>
    </Container>
  )
}

export default Analytics