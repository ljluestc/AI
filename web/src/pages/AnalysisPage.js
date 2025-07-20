import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { 
  Typography, 
  Box, 
  Paper, 
  Button, 
  CircularProgress, 
  Tabs, 
  Tab,
  Divider,
  Grid,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Chip,
  Accordion,
  AccordionSummary,
  AccordionDetails
} from '@mui/material';
import { 
  PlayArrow as AnalyzeIcon, 
  BugReport as BugIcon,
  Error as ErrorIcon,
  Warning as WarningIcon,
  Info as InfoIcon,
  ExpandMore as ExpandMoreIcon,
  ArrowBack as ArrowBackIcon,
  Code as CodeIcon
} from '@mui/icons-material';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/esm/styles/prism';
import axios from 'axios';

const GraphVisualization = ({ data }) => {
  // In a real app, this would use D3.js or a similar library to render a graph
  return (
    <Box sx={{ 
      height: 400, 
      bgcolor: 'rgba(255,255,255,0.03)', 
      borderRadius: 2, 
      p: 2, 
      display: 'flex', 
      justifyContent: 'center', 
      alignItems: 'center',
      border: '1px dashed rgba(255,255,255,0.2)'
    }}>
      <Typography variant="body1" color="text.secondary">
        Interactive Code Graph Visualization
      </Typography>
    </Box>
  );
};

const AnalysisPage = () => {
  const { repoId } = useParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [analyzing, setAnalyzing] = useState(false);
  const [repository, setRepository] = useState(null);
  const [analysisResults, setAnalysisResults] = useState(null);
  const [activeTab, setActiveTab] = useState(0);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchRepository = async () => {
      try {
        // In a real app, this would call the backend
        // const response = await axios.get(`/api/repositories/${repoId}`);
        
        // Simulate API call
        setTimeout(() => {
          // Extract repo name from path (repoId)
          const repoName = repoId.split('-').slice(0, -1).join('-');
          
          setRepository({
            id: repoId,
            name: repoName,
            path: repoId,
            clonedAt: new Date().toISOString(),
            lastScan: new Date(Date.now() - 86400000).toISOString(),
            errorsNum: 12
          });
          setLoading(false);
        }, 1000);
      } catch (err) {
        console.error('Error fetching repository:', err);
        setError('Failed to load repository');
        setLoading(false);
      }
    };

    fetchRepository();
  }, [repoId]);

  const handleStartAnalysis = async () => {
    setAnalyzing(true);
    try {
      // In a real app, this would call the backend
      // const response = await axios.post(`/api/analyze`, { repoPath: repoId });
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 3000));
      
      // Mock analysis results
      setAnalysisResults({
        analysisResults: [
          { file: 'src/components/auth/Login.js', line: 42, column: 5, message: 'Potential memory leak in useEffect hook', level: 'warning' },
          { file: 'src/services/api.js', line: 15, column: 3, message: 'API key exposed in client-side code', level: 'error' },
          { file: 'src/utils/helpers.js', line: 23, column: 10, message: 'Unused variable', level: 'info' },
          { file: 'src/components/Dashboard.js', line: 78, column: 15, message: 'React Hook useEffect has a missing dependency', level: 'warning' },
          { file: 'src/reducers/userReducer.js', line: 35, column: 7, message: 'Potential infinite loop', level: 'error' }
        ],
        errorLogs: [
          { file: 'src/services/api.js', line: 15, message: 'SecurityError: API key should not be exposed' },
          { file: 'src/reducers/userReducer.js', line: 35, message: 'Maximum update depth exceeded' }
        ],
        diagnosis: {
          'node-1': 'API key is exposed in the client-side code, causing security vulnerability',
          'node-2': 'Reducer function may update state in a way that triggers infinite re-renders'
        }
      });
      
      // Update repository with new scan time
      setRepository({
        ...repository,
        lastScan: new Date().toISOString(),
        errorsNum: 5
      });
    } catch (err) {
      console.error('Error during analysis:', err);
      setError('Analysis failed');
    } finally {
      setAnalyzing(false);
    }
  };

  const handleTabChange = (event, newValue) => {
    setActiveTab(newValue);
  };

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
        <Button 
          variant="contained" 
          sx={{ mt: 2 }} 
          startIcon={<ArrowBackIcon />}
          onClick={() => navigate('/repositories')}
        >
          Back to Repositories
        </Button>
      </Box>
    );
  }

  return (
    <Box sx={{ mb: 4 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <Button 
          startIcon={<ArrowBackIcon />} 
          onClick={() => navigate('/repositories')}
          sx={{ mr: 2 }}
        >
          Back
        </Button>
        <Typography variant="h4" sx={{ fontWeight: 500 }}>
          Repository Analysis: {repository.name}
        </Typography>
      </Box>

      <Paper sx={{ p: 3, mb: 3, border: '1px solid rgba(255,255,255,0.1)' }}>
        <Grid container spacing={2}>
          <Grid item xs={12} md={8}>
            <Typography variant="h6" gutterBottom>
              Repository Details
            </Typography>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
              <Typography variant="body1">
                <strong>Path:</strong> {repository.path}
              </Typography>
              <Typography variant="body1">
                <strong>Cloned:</strong> {new Date(repository.clonedAt).toLocaleString()}
              </Typography>
              {repository.lastScan && (
                <Typography variant="body1">
                  <strong>Last Analyzed:</strong> {new Date(repository.lastScan).toLocaleString()}
                </Typography>
              )}
              <Typography variant="body1">
                <strong>Status:</strong>{' '}
                {repository.errorsNum > 0 ? (
                  <Chip 
                    label={`${repository.errorsNum} issues found`} 
                    color="error" 
                    size="small" 
                    variant="outlined" 
                  />
                ) : (
                  <Chip 
                    label="No issues" 
                    color="success" 
                    size="small" 
                    variant="outlined" 
                  />
                )}
              </Typography>
            </Box>
          </Grid>
          <Grid item xs={12} md={4} sx={{ display: 'flex', alignItems: 'center', justifyContent: { xs: 'flex-start', md: 'flex-end' } }}>
            <Button
              variant="contained"
              startIcon={analyzing ? <CircularProgress size={20} color="inherit" /> : <AnalyzeIcon />}
              onClick={handleStartAnalysis}
              disabled={analyzing}
              sx={{ minWidth: 150 }}
            >
              {analyzing ? 'Analyzing...' : 'Start Analysis'}
            </Button>
          </Grid>
        </Grid>
      </Paper>

      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 3 }}>
        <Tabs value={activeTab} onChange={handleTabChange} aria-label="analysis tabs">
          <Tab label="Analysis Results" />
          <Tab label="Code Graph" />
          <Tab label="Error Logs" />
          <Tab label="Diagnosis" />
        </Tabs>
      </Box>

      {!analysisResults ? (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <CodeIcon sx={{ fontSize: 60, color: 'text.secondary', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No Analysis Results Yet
          </Typography>
          <Typography variant="body1" color="text.secondary" sx={{ mb: 3, maxWidth: 500, mx: 'auto' }}>
            Click the "Start Analysis" button to analyze this repository and get detailed insights about code quality and potential issues.
          </Typography>
          <Button
            variant="contained"
            startIcon={analyzing ? <CircularProgress size={20} color="inherit" /> : <AnalyzeIcon />}
            onClick={handleStartAnalysis}
            disabled={analyzing}
          >
            {analyzing ? 'Analyzing...' : 'Start Analysis'}
          </Button>
        </Box>
      ) : (
        <>
          {/* Analysis Results Tab */}
          {activeTab === 0 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Static Analysis Results
              </Typography>
              <List>
                {analysisResults.analysisResults.map((result, index) => (
                  <ListItem 
                    key={index}
                    sx={{ 
                      mb: 1, 
                      backgroundColor: 'rgba(255,255,255,0.03)',
                      borderRadius: 1,
                      border: '1px solid rgba(255,255,255,0.1)'
                    }}
                  >
                    <ListItemIcon>
                      {result.level === 'error' ? (
                        <ErrorIcon color="error" />
                      ) : result.level === 'warning' ? (
                        <WarningIcon color="warning" />
                      ) : (
                        <InfoIcon color="info" />
                      )}
                    </ListItemIcon>
                    <ListItemText
                      primary={result.message}
                      secondary={
                        <>
                          <Typography component="span" variant="body2" color="text.primary">
                            {result.file}
                          </Typography>
                          {` — Line ${result.line}, Column ${result.column}`}
                        </>
                      }
                    />
                    <Chip 
                      label={result.level} 
                      color={
                        result.level === 'error' 
                          ? 'error' 
                          : result.level === 'warning' 
                            ? 'warning' 
                            : 'info'
                      }
                      size="small"
                    />
                  </ListItem>
                ))}
              </List>
            </Box>
          )}

          {/* Code Graph Tab */}
          {activeTab === 1 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Code Knowledge Graph
              </Typography>
              <GraphVisualization data={{}} />
              <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
                The code knowledge graph visualizes relationships between different components in your code, 
                helping to identify dependencies and potential issues.
              </Typography>
            </Box>
          )}

          {/* Error Logs Tab */}
          {activeTab === 2 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Error Logs
              </Typography>
              <List>
                {analysisResults.errorLogs.map((log, index) => (
                  <ListItem 
                    key={index}
                    sx={{ 
                      mb: 1, 
                      backgroundColor: 'rgba(255,255,255,0.03)',
                      borderRadius: 1,
                      border: '1px solid rgba(255,255,255,0.1)'
                    }}
                  >
                    <ListItemIcon>
                      <BugIcon color="error" />
                    </ListItemIcon>
                    <ListItemText
                      primary={log.message}
                      secondary={`${log.file} — Line ${log.line}`}
                    />
                  </ListItem>
                ))}
              </List>
            </Box>
          )}

          {/* Diagnosis Tab */}
          {activeTab === 3 && (
            <Box>
              <Typography variant="h6" gutterBottom>
                Root Cause Diagnosis
              </Typography>
              {Object.entries(analysisResults.diagnosis).map(([nodeId, diagnosis], index) => (
                <Accordion 
                  key={nodeId}
                  sx={{ 
                    mb: 2, 
                    backgroundColor: 'rgba(255,255,255,0.03)',
                    border: '1px solid rgba(255,255,255,0.1)'
                  }}
                >
                  <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    aria-controls={`panel${index}-content`}
                    id={`panel${index}-header`}
                  >
                    <Typography sx={{ fontWeight: 500 }}>Issue {index + 1}: {diagnosis.split(':')[0]}</Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Typography gutterBottom>{diagnosis}</Typography>
                    <Divider sx={{ my: 2 }} />
                    <Typography variant="subtitle2" gutterBottom>Sample Code:</Typography>
                    <SyntaxHighlighter 
                      language="javascript" 
                      style={vscDarkPlus}
                      customStyle={{ borderRadius: 8 }}
                    >
                      {index === 0 
                        ? `// This is an example of the problematic code
const API_KEY = "1a2b3c4d5e6f7g8h9i0j"; // Security risk: API key exposed
                        
export const fetchData = async (endpoint) => {
  const response = await fetch(\`https://api.example.com/\${endpoint}?key=\${API_KEY}\`);
  return response.json();
};`
                        : `// This is an example of the problematic code
function userReducer(state, action) {
  switch (action.type) {
    case 'UPDATE_USER':
      // Problem: This may cause an infinite loop
      const newState = { ...state };
      newState.lastUpdated = new Date();
      dispatch({ type: 'LOG_UPDATE', data: newState }); // Dispatch inside reducer
      return newState;
    default:
      return state;
  }
}`
                      }
                    </SyntaxHighlighter>
                    <Typography variant="subtitle2" sx={{ mt: 2 }} gutterBottom>Recommended Fix:</Typography>
                    <SyntaxHighlighter 
                      language="javascript" 
                      style={vscDarkPlus}
                      customStyle={{ borderRadius: 8 }}
                    >
                      {index === 0 
                        ? `// Move API key to environment variables on the server
// Access via process.env in Node.js or use a backend API endpoint
                        
export const fetchData = async (endpoint) => {
  // The API key is now handled securely on the server
  const response = await fetch(\`https://api.example.com/proxy/\${endpoint}\`);
  return response.json();
};`
                        : `// Fixed version
function userReducer(state, action) {
  switch (action.type) {
    case 'UPDATE_USER':
      // Don't dispatch actions inside a reducer
      return {
        ...state,
        lastUpdated: new Date()
      };
    case 'LOG_UPDATE':
      // Handle logging in a separate case or middleware
      return state;
    default:
      return state;
  }
}`
                      }
                    </SyntaxHighlighter>
                  </AccordionDetails>
                </Accordion>
              ))}
            </Box>
          )}
        </>
      )}
    </Box>
  );
};

export default AnalysisPage;
