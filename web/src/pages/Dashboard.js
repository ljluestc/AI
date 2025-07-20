import React, { useState, useEffect } from 'react';
import { 
  Typography, 
  Grid, 
  Paper, 
  Box, 
  CircularProgress, 
  Button,
  Card,
  CardContent,
  CardActions,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListItemIcon
} from '@mui/material';
import { 
  BugReport as BugIcon, 
  Storage as StorageIcon, 
  CheckCircle as CheckIcon,
  Warning as WarningIcon,
  ArrowForward as ArrowForwardIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const Dashboard = () => {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [recentRepos, setRecentRepos] = useState([]);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get('/api/repositories');
        setRecentRepos(response.data.slice(0, 3)); // Get only the 3 most recent repos
        setLoading(false);
      } catch (err) {
        console.error('Error fetching repositories:', err);
        setError('Failed to load repositories. Please try again later.');
        setLoading(false);
      }
    };

    // Simulate API call with timeout
    setTimeout(() => {
      // Mock data for demo
      setRecentRepos([
        { id: 'repo-1', name: 'web-application', path: 'web-application-20231215121530', clonedAt: new Date().toISOString() },
        { id: 'repo-2', name: 'backend-api', path: 'backend-api-20231214093045', clonedAt: new Date(Date.now() - 86400000).toISOString() },
        { id: 'repo-3', name: 'mobile-app', path: 'mobile-app-20231213142210', clonedAt: new Date(Date.now() - 172800000).toISOString() }
      ]);
      setLoading(false);
    }, 1000);
  }, []);

  // Stats cards data
  const statsCards = [
    {
      title: 'Repositories',
      value: '12',
      icon: <StorageIcon sx={{ fontSize: 40, color: '#0984e3' }} />,
      color: '#0984e3',
    },
    {
      title: 'Issues Found',
      value: '248',
      icon: <BugIcon sx={{ fontSize: 40, color: '#d63031' }} />,
      color: '#d63031',
    },
    {
      title: 'Fixed Issues',
      value: '186',
      icon: <CheckIcon sx={{ fontSize: 40, color: '#00b894' }} />,
      color: '#00b894',
    }
  ];

  // Recent issues data
  const recentIssues = [
    {
      id: 'issue-1',
      repo: 'web-application',
      file: 'src/components/auth/Login.js',
      message: 'Potential memory leak in useEffect hook',
      severity: 'warning',
    },
    {
      id: 'issue-2',
      repo: 'backend-api',
      file: 'api/controllers/user.go',
      message: 'SQL injection vulnerability in query parameter',
      severity: 'error',
    },
    {
      id: 'issue-3',
      repo: 'mobile-app',
      file: 'app/screens/Profile.tsx',
      message: 'Unused variable declaration',
      severity: 'info',
    },
  ];

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return (
      <Box sx={{ textAlign: 'center', mt: 5 }}>
        <Typography color="error" variant="h6">{error}</Typography>
        <Button variant="contained" sx={{ mt: 2 }} onClick={() => window.location.reload()}>
          Try Again
        </Button>
      </Box>
    );
  }

  return (
    <Box sx={{ mb: 4 }}>
      <Typography variant="h4" gutterBottom sx={{ fontWeight: 500 }}>
        Dashboard
      </Typography>
      
      {/* Stats Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {statsCards.map((card) => (
          <Grid item xs={12} md={4} key={card.title}>
            <Paper
              elevation={2}
              sx={{
                p: 3,
                display: 'flex',
                alignItems: 'center',
                height: '100%',
                background: `linear-gradient(135deg, rgba(30,30,30,1) 0%, rgba(40,40,40,1) 100%)`,
                border: '1px solid rgba(255,255,255,0.1)',
              }}
            >
              <Box sx={{ mr: 3 }}>{card.icon}</Box>
              <Box>
                <Typography variant="h6" color="text.secondary">{card.title}</Typography>
                <Typography variant="h3" sx={{ fontWeight: 600, color: card.color }}>{card.value}</Typography>
              </Box>
            </Paper>
          </Grid>
        ))}
      </Grid>

      <Grid container spacing={3}>
        {/* Recent Repositories */}
        <Grid item xs={12} md={6}>
          <Paper 
            elevation={2} 
            sx={{ 
              p: 2, 
              height: '100%', 
              border: '1px solid rgba(255,255,255,0.1)' 
            }}
          >
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6">Recent Repositories</Typography>
              <Button 
                size="small" 
                endIcon={<ArrowForwardIcon />}
                onClick={() => navigate('/repositories')}
              >
                View All
              </Button>
            </Box>
            <Divider sx={{ mb: 2 }} />
            
            {recentRepos.length > 0 ? (
              <Grid container spacing={2}>
                {recentRepos.map((repo) => (
                  <Grid item xs={12} key={repo.id}>
                    <Card variant="outlined" sx={{ backgroundColor: 'rgba(255,255,255,0.03)' }}>
                      <CardContent sx={{ pb: 1 }}>
                        <Typography variant="h6" component="div">
                          {repo.name}
                        </Typography>
                        <Typography color="text.secondary" sx={{ fontSize: 14 }}>
                          Last updated: {new Date(repo.clonedAt).toLocaleDateString()}
                        </Typography>
                      </CardContent>
                      <CardActions>
                        <Button 
                          size="small" 
                          onClick={() => navigate(`/analysis/${repo.path}`)}
                        >
                          Analyze
                        </Button>
                        <Button size="small">View Code</Button>
                      </CardActions>
                    </Card>
                  </Grid>
                ))}
              </Grid>
            ) : (
              <Box sx={{ textAlign: 'center', py: 4 }}>
                <Typography color="text.secondary">No repositories found</Typography>
                <Button 
                  variant="contained" 
                  sx={{ mt: 2 }}
                  onClick={() => navigate('/repositories')}
                >
                  Add Repository
                </Button>
              </Box>
            )}
          </Paper>
        </Grid>

        {/* Recent Issues */}
        <Grid item xs={12} md={6}>
          <Paper 
            elevation={2} 
            sx={{ 
              p: 2, 
              height: '100%',
              border: '1px solid rgba(255,255,255,0.1)' 
            }}
          >
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6">Recent Issues</Typography>
              <Button 
                size="small" 
                endIcon={<ArrowForwardIcon />}
                onClick={() => navigate('/errors')}
              >
                View All
              </Button>
            </Box>
            <Divider sx={{ mb: 2 }} />
            
            <List>
              {recentIssues.map((issue) => (
                <ListItem 
                  key={issue.id}
                  sx={{ 
                    mb: 1, 
                    backgroundColor: 'rgba(255,255,255,0.03)',
                    borderRadius: 1,
                    '&:hover': {
                      backgroundColor: 'rgba(255,255,255,0.05)'
                    }
                  }}
                >
                  <ListItemIcon>
                    {issue.severity === 'error' ? (
                      <BugIcon color="error" />
                    ) : issue.severity === 'warning' ? (
                      <WarningIcon color="warning" />
                    ) : (
                      <CheckIcon color="info" />
                    )}
                  </ListItemIcon>
                  <ListItemText
                    primary={issue.message}
                    secondary={
                      <>
                        <Typography component="span" variant="body2" color="text.primary">
                          {issue.repo}
                        </Typography>
                        {` — ${issue.file}`}
                      </>
                    }
                  />
                </ListItem>
              ))}
            </List>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;
