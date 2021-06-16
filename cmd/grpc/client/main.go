package main

import (
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/client_to_server/song_client_service"
	"google.golang.org/grpc"
	"log"
)


func authMethods() map[string]bool {
	const roleServicePath = "/role.RoleService/"
	const userServicePath = "/user.UserService/"
	const subscriptionServicePath = "/subscription.SubscriptionService/"
	const subscriptionTypeServicePath = "/subscription.SubscriptionTypeService/"
	const bookingServicePath = "/booking.BookingService/"
	const addressServicePath = "/address.AddressService/"
	const songServicePath = "/song.SongService/"

	return map[string]bool{
		roleServicePath + "ListRoles":                          true,
		roleServicePath + "GetRole":                            true,
		roleServicePath + "CreateRole":                         true,
		roleServicePath + "DeleteRole":                         true,
		roleServicePath + "UpdateRole":                         true,
		userServicePath + "ListUsers":                          true,
		userServicePath + "GetUser":                            true,
		userServicePath + "CreateUser":                         true,
		userServicePath + "DeleteUser":                         true,
		userServicePath + "UpdateUser":                         true,
		userServicePath + "Login":                              true,
		userServicePath + "SignUp":                             true,
		subscriptionServicePath + "ListSubscriptions":          true,
		subscriptionServicePath + "GetSubscription":            true,
		subscriptionServicePath + "CreateSubscription":         true,
		subscriptionServicePath + "DeleteSubscription":         true,
		subscriptionServicePath + "UpdateSubscription":         true,
		subscriptionTypeServicePath + "ListSubscriptionTypes":  true,
		subscriptionTypeServicePath + "GetSubscriptionType":    true,
		subscriptionTypeServicePath + "CreateSubscriptionType": true,
		subscriptionTypeServicePath + "DeleteSubscriptionType": true,
		subscriptionTypeServicePath + "UpdateSubscriptionType": true,
		bookingServicePath + "ListBookings":                    true,
		bookingServicePath + "GetBooking":                      true,
		bookingServicePath + "CreateBooking":                   true,
		bookingServicePath + "DeleteBooking":                   true,
		bookingServicePath + "UpdateBookings":                  true,
		addressServicePath + "ListAddresses":                   true,
		addressServicePath + "GetAddress":                      true,
		addressServicePath + "CreateAddress":                   true,
		addressServicePath + "DeleteAddress":                   true,
		addressServicePath + "UpdateAddress":                   true,
		songServicePath + "CreateSong":                         true,
		songServicePath + "GetAllSongs":                  		true,
		songServicePath + "GetSong":                  			true,
		songServicePath + "UpdateSong":                  		true,
		songServicePath + "DeleteSong":                  		true,
		songServicePath + "LikeSong":                  			true,
		songServicePath + "GetLikesCountRequest":               true,
		songServicePath + "GetSongViewCount":                  	true,
		songServicePath + "IncreaseAlbumViewCount":             true,
	}
}
func main() {
	transportOption := grpc.WithInsecure()
	grpcDialStr, err := constant.GetGrpcConnectionString()
	if err != nil {
		log.Fatal("Cannot dial to the server")
	}

	//interceptor, err := auth_client_service.NewAuthInterceptor(authMethods(), "dsf")
	//if err != nil {
	//	log.Fatal("cannot create auth interceptor: ", err)
	//}
	cc2, err := grpc.Dial(
		grpcDialStr,
		transportOption,
		//grpc.WithUnaryInterceptor(interceptor.Unary()),
		//grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	songClient:= song_client_service.NewSongClient(cc2)
	TestGetSong(songClient)
	TestGetAllSongs(songClient)
	TestCreateSong(songClient)

}



