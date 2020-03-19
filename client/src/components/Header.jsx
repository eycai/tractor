import React from 'react';
import { Pane, Text } from 'evergreen-ui';

const Header = () => (
  <Pane
    height={60}
    display="flex"
    background="#084B8A"
    alignItems="center"
    padding={24}
  >
    <Text color="white">
      <h1>Tractor</h1>
    </Text>
  </Pane>
);

export default Header;
