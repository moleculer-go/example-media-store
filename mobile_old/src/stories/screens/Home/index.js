import * as React from "react";
import {
  Container,
  Header,
  Title,
  Content,
  Text,
  Button,
  Icon,
  Left,
  Body,
  Right,
  List,
  ListItem
} from "native-base";
import { Col, Row, Grid } from "react-native-easy-grid";

import styles from "./styles";
export interface Props {
  navigation: any;
  list: any;
}
export interface State {}
class Home extends React.Component<Props, State> {
  render() {
    return (
      <Container style={styles.container}>
        <Header>
          <Left>
            <Button transparent>
              <Icon
                active
                name="menu"
                onPress={() => this.props.navigation.navigate("DrawerOpen")}
              />
            </Button>
          </Left>

          <Body>
            <Title>Picture Share</Title>
          </Body>

          <Right />
        </Header>
        <Content>
        <Grid>
            <Col>
                <Text>1</Text>
            </Col>
            <Col>
                <Row>
                    <Text>2</Text>
                </Row>
                <Row>
                    <Text>3</Text>
                </Row>
            </Col>
        </Grid>
          {/* <List>
            {this.props.list.map((item, i) => (
              <ListItem
                key={i}
                onPress={() =>
                  this.props.navigation.navigate("BlankPage", {
                    name: { item }
                  })}
              >
                <Text>{item}</Text>
              </ListItem>
            ))}
          </List> */}
        </Content>
      </Container>
    );
  }
}

export default Home;
