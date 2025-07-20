import React, { useState, useEffect } from 'react';
import { 
  Typography, 
  Paper, 
  Button, 
  TextField, 
  Dialog, 
  DialogActions, 
  DialogContent, 
  DialogContentText, 
  DialogTitle, 
  Table, 
  TableBody, 
  TableCell, 
  TableContainer, 
  TableHead, 
  TableRow, 
  Box, 
  CircularProgress,
  IconButton,
  Tooltip,
  Chip
} from '@mui/material';
import { 
  Add as AddIcon, 
  Refresh as RefreshIcon, 
  PlayArrow as AnalyzeIcon, 
  Visibility as ViewIcon,
  DeleteOutline as DeleteIcon
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const Repositories = () => {
  const navigate = useNavigate();
  const [repositories, setRepositories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [openDialog, setOpenDialog] = useState(false);
  const [repoUrl, setRepoUrl] = useState('');
  const [cloning, setCloning] = useState(false);

  useEffect(() => {
    fetchRepositories();
  }, []);

  const fetchRepositories = async () => {
    setLoading(true);
    try {
      // Simulate API call
      setTimeout(() => {
        // Mock data for demo
        setRepositories([
          { id: 'repo-1', name: 'web-application', path: 'web-application-20231215121530', clonedAt: new Date().toISOString(), errorsNum: 24 },
          { id: 'repo-2', name: 'backend-api', path: 'backend-api-20231214093045', clonedAt: new Date(Date.now() - 86400000).toISOString(), errorsNum: 15 },
          { id: 'repo-3', name: 'mobile-app', path: 'mobile-app-20231213142210', clonedAt: new Date(Date.now() - 172800000).toISOString(), errorsNum: 7 },
          { id: 'repo-4', name: 'data-processor', path: 'data-processor-20231212105530', clonedAt: new Date(Date.now() - 259200000).toISOString(), errorsNum: 3 },
          { id: 'repo-5', name: 'authentication-service', path: 'authentication-service-20231211133045', clonedAt: new Date(Date.now() - 345600000).toISOString(), errorsNum: 0 }
        ]);
        setLoading(false);
      }, 1000);
    } catch (err) {
      console.error('Error fetching repositories:', err);
      setError('Failed to load repositories');
      setLoading(false);
    }
  };

  const handleOpenDialog = () => {
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
    setRepoUrl('');
  };

  const handleCloneRepo = async () => {
    if (!repoUrl) return;
    
    setCloning(true);
    try {
      // In a real app, this would call the backend
      // const response = await axios.post('/api/repositories', { url: repoUrl });
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 2000));
      
      // Add mock repository to the list
      const repoName = repoUrl.split('/').pop().replace('.git', '');
      const timestamp = new Date().toISOString().replace(/[-:]/g, '').substring(0, 14);
      const newRepo = {
        id: `repo-${repositories.length + 1}`,
        name: repoName,
        path: `${repoName}-${timestamp}`,
        clonedAt: new Date().toISOString(),
        errorsNum: 0
      };
      
      setRepositories([newRepo, ...repositories]);
      handleCloseDialog();
    } catch (err) {
      console.error('Error cloning repository:', err);
      setError('Failed to clone repository');
    } finally {
      setCloning(false);
    }
  };

  const handleAnalyzeRepo = (repoPath) => {
    navigate(`/analysis/${repoPath}`);
  };

  if (loading && repositories.length === 0) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '80vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box sx={{ mb: 4 }}>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4" gutterBottom sx={{ fontWeight: 500 }}>
          Repositories
        </Typography>
        <Box>
          <Button 
            variant="outlined" 
            startIcon={<RefreshIcon />} 
            onClick={fetchRepositories}
            sx={{ mr: 2 }}
          >
            Refresh
          </Button>
          <Button 
            variant="contained" 
            startIcon={<AddIcon />} 
            onClick={handleOpenDialog}
          >
            Add Repository
          </Button>
        </Box>
      </Box>

      {error && (
        <Paper elevation={0} sx={{ p: 2, mb: 3, bgcolor: 'error.dark', color: 'white' }}>
          <Typography>{error}</Typography>
        </Paper>
      )}

      <TableContainer component={Paper} sx={{ border: '1px solid rgba(255,255,255,0.1)' }}>
        <Table sx={{ minWidth: 650 }} aria-label="repositories table">
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Path</TableCell>
              <TableCell>Cloned At</TableCell>
              <TableCell>Status</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {repositories.map((repo) => (
              <TableRow
                key={repo.id}
                sx={{ '&:last-child td, &:last-child th': { border: 0 }, '&:hover': { backgroundColor: 'rgba(255,255,255,0.03)' } }}
              >
                <TableCell component="th" scope="row">
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {repo.name}
                  </Typography>
                </TableCell>
                <TableCell>{repo.path}</TableCell>
                <TableCell>{new Date(repo.clonedAt).toLocaleString()}</TableCell>
                <TableCell>
                  {repo.errorsNum > 0 ? (
                    <Chip 
                      label={`${repo.errorsNum} issues found`} 
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
                </TableCell>
                <TableCell align="right">
                  <Tooltip title="Analyze">
                    <IconButton 
                      color="primary" 
                      onClick={() => handleAnalyzeRepo(repo.path)}
                    >
                      <AnalyzeIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="View Code">
                    <IconButton color="secondary">
                      <ViewIcon />
                    </IconButton>
                  </Tooltip>
                  <Tooltip title="Delete">
                    <IconButton color="error">
                      <DeleteIcon />
                    </IconButton>
                  </Tooltip>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      {/* Add Repository Dialog */}
      <Dialog open={openDialog} onClose={handleCloseDialog}>
        <DialogTitle>Add New Repository</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Enter the Git URL of the repository you want to analyze.
          </DialogContentText>
          <TextField
            autoFocus
            margin="dense"
            id="repoUrl"
            label="Repository URL"
            type="text"
            fullWidth
            variant="outlined"
            value={repoUrl}
            onChange={(e) => setRepoUrl(e.target.value)}
            placeholder="https://github.com/username/repository.git"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog} disabled={cloning}>
            Cancel
          </Button>
          <Button 
            onClick={handleCloneRepo} 
            variant="contained" 
            disabled={!repoUrl || cloning}
            startIcon={cloning ? <CircularProgress size={20} /> : null}
          >
            {cloning ? 'Cloning...' : 'Clone Repository'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default Repositories;
