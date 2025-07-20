import React from 'react';
import { 
  Drawer, 
  List, 
  ListItem, 
  ListItemIcon, 
  ListItemText, 
  Divider, 
  Box, 
  Toolbar,
  Typography,
  ListItemButton
} from '@mui/material';
import { 
  Dashboard as DashboardIcon, 
  Storage as StorageIcon,
  BarChart as BarChartIcon,
  BugReport as BugReportIcon,
  Code as CodeIcon,
  Settings as SettingsIcon,
  Help as HelpIcon
} from '@mui/icons-material';
import { useNavigate, useLocation } from 'react-router-dom';

const drawerWidth = 240;

const Sidebar = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    {
      text: 'Dashboard',
      icon: <DashboardIcon />,
      path: '/'
    },
    {
      text: 'Repositories',
      icon: <StorageIcon />,
      path: '/repositories'
    },
    {
      text: 'Analysis',
      icon: <BarChartIcon />,
      path: '/analysis'
    },
    {
      text: 'Error Tracking',
      icon: <BugReportIcon />,
      path: '/errors'
    },
    {
      text: 'Code Explorer',
      icon: <CodeIcon />,
      path: '/code'
    }
  ];

  const supportItems = [
    {
      text: 'Settings',
      icon: <SettingsIcon />,
      path: '/settings'
    },
    {
      text: 'Help & Support',
      icon: <HelpIcon />,
      path: '/help'
    }
  ];

  return (
    <Drawer
      variant="permanent"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        [`& .MuiDrawer-paper`]: { 
          width: drawerWidth, 
          boxSizing: 'border-box',
          background: (theme) => theme.palette.background.paper,
          borderRight: '1px solid rgba(255, 255, 255, 0.12)'
        },
      }}
    >
      <Toolbar />
      <Box sx={{ overflow: 'auto', pt: 2 }}>
        <List>
          {menuItems.map((item) => (
            <ListItem key={item.text} disablePadding>
              <ListItemButton
                selected={location.pathname === item.path}
                onClick={() => navigate(item.path)}
                sx={{
                  '&.Mui-selected': {
                    backgroundColor: 'rgba(0, 184, 148, 0.1)',
                    '&:hover': {
                      backgroundColor: 'rgba(0, 184, 148, 0.2)',
                    },
                  },
                  '&:hover': {
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                  },
                  borderLeft: location.pathname === item.path ? '4px solid #00b894' : '4px solid transparent',
                  pl: 2,
                }}
              >
                <ListItemIcon sx={{ color: location.pathname === item.path ? '#00b894' : 'inherit' }}>
                  {item.icon}
                </ListItemIcon>
                <ListItemText 
                  primary={item.text} 
                  primaryTypographyProps={{ 
                    fontWeight: location.pathname === item.path ? 600 : 400 
                  }}
                />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
        <Divider sx={{ my: 2, backgroundColor: 'rgba(255, 255, 255, 0.12)' }} />
        <Typography variant="overline" sx={{ px: 3, color: 'text.secondary' }}>
          Support
        </Typography>
        <List>
          {supportItems.map((item) => (
            <ListItem key={item.text} disablePadding>
              <ListItemButton
                onClick={() => navigate(item.path)}
                sx={{
                  '&:hover': {
                    backgroundColor: 'rgba(255, 255, 255, 0.05)',
                  },
                  pl: 2,
                }}
              >
                <ListItemIcon>{item.icon}</ListItemIcon>
                <ListItemText primary={item.text} />
              </ListItemButton>
            </ListItem>
          ))}
        </List>
      </Box>
    </Drawer>
  );
};

export default Sidebar;
