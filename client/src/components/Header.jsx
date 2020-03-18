import React from 'react';
import { Pane, Text } from 'evergreen-ui';

const Header = () => (
  <Pane
    height={60}
    display="flex"
    background="tealTint"
    alignItems="center"
    padding={24}
    marginBottom={16}
  >
    <Text>
      <h1>Tractor</h1>
    </Text>
  </Pane>
);

export default Header;
