    <template>
      <div></div>
    </template>
    
    <script>
    import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core';
    import { setContext } from '@apollo/client/link/context';
    import gql from 'graphql-tag';
    
    export default {
      name: 'GraphQLTester',
      data() {
        return {
          rooms: [],
          pageInfo: {},
          filteredRooms: [],
          startOfDayTimestamp: null,
          endOfDayTimestamp: null,
          ids: [],
          edges: [],
          users: [],
          eventList: [],
        };
      },
      created() {
        const httpLink = createHttpLink({
          uri: 'http://localhost:8080/query', 
        });

        const authLink = setContext((_, { headers }) => {
          return {
            headers: {
              ...headers,
              authorization: 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImxlZWl2YW4xMDA3QGdtYWlsLmNvbSIsImV4cCI6MTcxODA2ODI5OCwibmFtZSI6Ikl2YW4gTGVlIiwic3ViIjoiNjY0NWVjZTEzNmUyYTBmMDM1OTYxYmRkIn0.u0a949cBKw2qy3uVOXikTTDGHiU5UN5eUROnpA5QHTw',
            }
          }
        });

        this.client = new ApolloClient({
          link: authLink.concat(httpLink),
          cache: new InMemoryCache(),
        });
      },
      methods: {
        createRoom(roomInput) {

          const CREATE_ROOM_MUTATION = gql`
            mutation myCreate($my_input: UpsertRoomInput!) {
              upsertRoom(room: $my_input) {
                id
                name
                capacity
                equipments
                rules
                isDelete
              }
            }
          `;

          const variables = {
            my_input: roomInput
          };

          return this.client.mutate({
            mutation: CREATE_ROOM_MUTATION,
            variables
          }).then(response => {
            console.log("Room created or updated successfully:", response.data);
          }).catch(error => {
            console.error("Error creating or updating room:", error);
          });
        },
        createEvent(eventInput) {

          const CREATE_EVENT_MUTATION = gql`
            mutation createEvent($input: UpsertEventInput!) {
              upsertEvent(input: $input) {
                id
                title
                description
                startAt
                endAt
                roomReservation {
                  room {
                    id
                    name
                    capacity
                    equipments
                    rules
                    isDelete
                  }
                  status
                }
                participants {
                  id
                  name
                  email
                }
                notes
                remindAt
                summary
                creator {
                  id
                  name
                  email
                }
                isDelete
                attachedFile {
                  url
                  name
                }
              }
            }
          `;

          const variables = {
            input: eventInput
          };
          return this.client.mutate({
            mutation: CREATE_EVENT_MUTATION,
            variables
          }).then(response => {
            console.log("Event created or updated successfully:", response.data);
          }).catch(error => {
            console.error("Error creating or updating room:", error);
          });
        },
        queryAllRooms() {

          const GET_ALL_ROOMS_QUERY = gql`
            query {
              paginatedRooms(first: 3, after: "") {
                edges {
                  node {
                    id
                    name
                    capacity
                    equipments
                    rules
                    isDelete
                  }
                }
                pageInfo {
                  startCursor
                  endCursor
                }
              }
            }
          `;

          this.client.query({
            query: GET_ALL_ROOMS_QUERY,
            fetchPolicy: 'no-cache'
          }).then(response => {
            const paginatedRooms = response.data.paginatedRooms;

            if (paginatedRooms && paginatedRooms.edges.length > 0) {
              this.rooms = paginatedRooms.edges.map(edge => edge.node);
              this.pageInfo = paginatedRooms.pageInfo;
              this.$emit('queryAllRooms', this.rooms);
            } else {
              this.$emit('queryAllRooms', []);
            }
          }).catch(error => {
            console.error("Failed to fetch rooms:", error);
          });

          
        },
        deleteRoom(roomId) {
          const DELETE_ROOM_MUTATION = gql`
            mutation deleteRoom($id: ID!) {
              deleteRoom(id: $id) {
                id
              }
            }
          `;
          return this.client.mutate({
            mutation: DELETE_ROOM_MUTATION,
            variables: {
              id: roomId
            }
          }).then(response => {
            console.log('Room deleted:', response.data.deleteRoom.id);
          }).catch(error => {
            console.error('Error deleting room:', error);
          });
        },
        deleteEvent(eventId) {
          const DELETE_EVENT_MUTATION = gql`
            mutation deleteEvent($id: ID!) {
              deleteEvent(id: $id) {
                roomReservation {
                  status
                }
                isDelete
              }
            }
          `;
          return this.client.mutate({
            mutation: DELETE_EVENT_MUTATION,
            variables: {
              id: eventId
            },
            fetchPolicy: 'no-cache',
          }).then(response => {
            console.log('Event deleted:', response.data.deleteEvent.isDelete);
            console.log('Room reservation status:', response.data.deleteEvent.roomReservation.status);
          }).catch(error => {
            console.error('Error deleting event:', error);
          });
        },
        fetchAvailableRooms(variables) {
        return new Promise((resolve, reject) => {

          this.calculateStartOfDay(variables.startAt);
          this.calculateEndOfDay(variables.endAt);

          // First queryRoomSchedules
          if (!variables.ids) {
            variables.ids = [];
          }

          
          this.queryRoomSchedules(variables)
            .then(response => {
              if (response.data.paginatedRoomSchedules !== null) {
                
                const edges = response.data.paginatedRoomSchedules.edges;
                this.rooms = edges.map(edge => {
                  return {
                    ...edge.node.room,
                    schedules: edge.node.schedules
                  };
                });

                // Second queryRoomSchedules inside the then block of the first
                variables.ids = this.rooms.map(room => room.id);
                variables.startAt = this.startOfDayTimestamp;
                variables.endAt = this.endOfDayTimestamp;

                return this.queryRoomSchedules(variables);
              } else {
                resolve([]); // 返回一个空数组
              }
            })
            .then(response => {
              const edges = response.data.paginatedRoomSchedules.edges;
              this.rooms = edges.map(edge => {
                return {
                  ...edge.node.room,
                  schedules: edge.node.schedules
                };
              });
              this.$emit('fetchAvailableRooms', this.rooms);
              resolve(this.rooms);
            })
            .catch(error => {
              this.error = error;
              // console.error("Error in queryRoomSchedules:", error);
              reject(error);
            });
          });
        },   
        editSummary(variables) {
          const UPDATE_SUMMARY_MUTATION = gql`
            mutation updateEventSummary($id: ID!, $summary: String!) {
              updateEventSummary(id: $id, summary: $summary)
            }
          `;
          return this.client.mutate({
            mutation: UPDATE_SUMMARY_MUTATION,
            variables
          }).then(response => {
            console.log('Summary edited and New summary', response);
          }).catch(error => {
            console.error('Error editing summary:', error);
          });
        },
        queryRoomSchedules(variables) {

          const GET_ROOM_SCHEDULES = gql`
            query getRoomSchedules(
              $ids: [String!]!,
              $startAt: Int64!,
              $endAt: Int64!,
              $rules: [Rule!],
              $equipments: [Equipment!],
              $first: Int = 20,
              $after: String
            ) {
              paginatedRoomSchedules(
                ids: $ids,
                startAt: $startAt,
                endAt: $endAt,
                rules: $rules,
                equipments: $equipments,
                first: $first,
                after: $after
              ) {
                edges {
                  node {
                    room {
                      id
                      name
                      capacity
                      equipments
                      rules
                      isDelete
                    }
                    schedules {
                      id
                      title
                      description
                      startAt
                      endAt
                      participants {
                        id
                        name
                      }
                      summary
                      creator {
                        id
                      }
                      attachedFile {
                        url
                        name
                      }
                    }
                  }
                  cursor
                }
                pageInfo {
                  endCursor
                }
              }
            }
          `;

          return this.client.query({
            query: GET_ROOM_SCHEDULES,
            variables, //: defaultVariables,
            fetchPolicy: 'no-cache',
          }).then(response => {
            return response
          }).catch(error => {
            this.error = error;
            console.error("Failed to fetch room schedules:", error);
          });
        },
        calculateStartOfDay(timestamp) {
          const date = new Date(timestamp);
          const year = date.getFullYear();
          const month = date.getMonth();
          const day = date.getDate();

          const startOfDay = new Date(year, month, day, 0, 0, 0, 0);

          this.startOfDayTimestamp = startOfDay.getTime();
        },
        calculateEndOfDay(timestamp) {
          const date = new Date(timestamp);
          const year = date.getFullYear();
          const month = date.getMonth();
          const day = date.getDate();

          const endOfDay = new Date(year, month, day, 23, 59, 59, 999);

          this.endOfDayTimestamp = endOfDay.getTime();
        },
        queryUsers() {

          const GET_PAGINATED_USERS = gql`
            query getPaginatedUsers($first: Int = 20, $after: String) {
              paginatedUsers(first: $first, after: $after) {
                edges {
                  node {
                    id
                    name
                  }
                }
              }
            }
          `;

          const variables = {
            first: 20,
            after: null
          };

          return this.client.query({
            query: GET_PAGINATED_USERS,
            variables
          }).then(response => {
            this.users = response.data.paginatedUsers.edges.map(edge => edge.node);
            this.$emit('queryUsers', this.users);
            return this.users; // 返回查詢到的用戶
          }).catch(error => {
            console.error('Error fetching users:', error);
            throw error; // 繼續拋出錯誤，以便在外部捕獲
          });
        },
        getUserEvents(variables) {
          const GET_USER_EVENTS = gql`
            query getUserEvents($userIDs: [ID!]!, $startAt: Int64!, $endAt: Int64!) {
              userEvents(userIDs: $userIDs, startAt: $startAt, endAt: $endAt) {
                events {
                  id
                  title
                  description
                  startAt
                  endAt
                  roomReservation {
                    room {
                      id
                      name
                      capacity
                      equipments
                      rules
                    }
                    status
                  }
                    participants {
                      id
                      name
                    }
                  summary
                  attachedFile {
                    url
                    name
                  }
                  creator {
                    id
                    name
                  }
                }
              }
            }
          `;
          this.client.query({
            query: GET_USER_EVENTS,
            variables,
            fetchPolicy: 'no-cache'
          }).then(response => {
            let returnEventList = [];
            if (response.data.userEvents.length > 0) {
              const events = response.data.userEvents[1].events;
              this.eventList = [];
              for (let i = 0; i < events.length; i++) {
                const event = events[i];
                const processedEvent = {
                  eventId: event.id,
                  title: event.title,
                  description: event.description,
                  fileName: event.attachedFile.name,
                  fileUrl: event.attachedFile.url,
                  creatorId: event.creator.id,
                  participants: event.participants,
                  startAt: event.startAt,
                  endAt: event.endAt,
                  summary: event.summary,
                  roomId: [event.roomReservation.room.id],
                };
                this.eventList.push(processedEvent);
              }
              
              console.log('this.eventList:', this.eventList);
              // Use the first event from the eventList to call fetchAvailableRooms
              if (this.eventList.length > 0) {
                let promises = this.eventList.map(event => {
                  const fetchVariables = {
                    startAt: event.startAt,
                    endAt: event.endAt,
                    ids: event.roomId
                  };
                  console.log('fetchVariables:', fetchVariables);
                  return this.fetchAvailableRooms(fetchVariables)
                    .then(rooms => {
                    return {
                      originalData: event,
                      roomsData: rooms
                    };
                    })
                    .catch(error => {
                      console.error('Error fetching available rooms:', error);
                      return {
                        originalData: event,
                        roomsData: null,
                        error: error
                      };
                    });
                });


                // Wait for all fetchAvailableRooms promises to resolve
                Promise.all(promises)
                  .then(results => {
                    returnEventList = results;
                    this.$emit('getEventList', returnEventList);
                  })
                  .catch(error => {
                    console.error('Error processing events:', error);
                  });

              }
            } else {

              this.$emit('getEventList', returnEventList);
            }
          }).catch(error => {
            console.error('Error fetching user events:', error);
          });
        }
      },
    }
    </script>
    