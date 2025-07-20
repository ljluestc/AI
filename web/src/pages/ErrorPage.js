import React from 'react';
import { Box, Typography, Button, Paper } from '@mui/material';
import { Error as ErrorIcon, Home as HomeIcon } from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';

const ErrorPage = () => {
  const navigate = useNavigate();

  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: '70vh',
        textAlign: 'center',
        p: 3
      }}
    >
      <Paper
        elevation={2}
        sx={{
          p: 5,
          borderRadius: 3,
          maxWidth: 500,
          background: 'linear-gradient(135deg, rgba(30,30,30,1) 0%, rgba(40,40,40,1) 100%)',
          border: '1px solid rgba(255,255,255,0.1)'
        }}
      >
        <ErrorIcon color="error" sx={{ fontSize: 80, mb: 2 }} />
        <Typography variant="h4" gutterBottom sx={{ fontWeight: 600 }}>
          Page Not Found
        </Typography>
        <Typography variant="body1" sx={{ mb: 3, color: 'text.secondary' }}>
          Sorry, the page you are looking for does not exist or has been moved.
        </Typography>
        <Button 
          variant="contained"
          startIcon={<HomeIcon />}
          onClick={() => navigate('/')}
          sx={{ minWidth: 150 }}
        >
          Back to Home
        </Button>
      </Paper>
    </Box>
  );
};

export default ErrorPage;
