import React from 'react';
import ReactDOM from 'react-dom/client';
import { ChakraProvider, extendTheme } from '@chakra-ui/react'

import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';

const globalFontSize = '14px';
const theme = extendTheme({
  styles: {
    global: {
      body: {
        fontSize: globalFontSize,
      },
    },
  },
  components: {
    Table: {
      baseStyle: {
        th: {
          fontSize: globalFontSize,
          borderBottom: '2px solid',
          borderColor: 'gray.200',
          bg: 'gray.50',
          fontWeight: 'bold',
          textTransform: 'uppercase',
          letterSpacing: 'wider',
        },
        td: {
          fontSize: globalFontSize,
          borderBottom: '1px solid',
          borderColor: 'gray.100',
        },
        table: {
          borderCollapse: 'collapse',
        },
        caption: {
          textAlign: 'left',
          fontSize: 'xl',
          fontWeight: 'bold',
        },
      },
    },
    Button: {
      baseStyle: {
        fontSize: globalFontSize,
      },
    },
  },
});

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <ChakraProvider theme={ theme }>
      <App />
    </ChakraProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
