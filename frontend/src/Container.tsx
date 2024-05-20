import React from 'react';

interface ContainerProps {
  title: string;
  children: React.ReactNode;
}

const Container: React.FC<ContainerProps> = ({ title, children }) => {
  return (
    <div style={styles.container}>
      <h2>{title}</h2>
      {children}
    </div>
  );
};

const styles = {
  container: {
    border: '1px solid #ccc',
    padding: '20px',
    borderRadius: '5px',
    backgroundColor: '#f9f9f9',
  } as React.CSSProperties,
};

export default Container;
