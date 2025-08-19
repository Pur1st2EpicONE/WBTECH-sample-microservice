package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/broker"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/configs"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

func main() {

	config, err := configs.ProdConfig()
	if err != nil {
		logger.LogFatal("producer — failed to load config", err)
	}

	producer, err := broker.NewProducer(config)
	if err != nil {
		logger.LogFatal("producer — creation failed", err)
	}

	checkArgs(&config.TotalMessages)
	orders := getOrders(config.TotalMessages)

	for i, order := range orders {
		orderJSON, err := json.MarshalIndent(order, "", "   ")
		if err != nil {
			logger.LogFatal("producer — failed to marshal order with indent", err)
		}
		logger.LogInfo(fmt.Sprintf("order-producer — sending order %d to Kafka", i+1))
		producer.Produce(orderJSON, config.Topic)
	}
}

func checkArgs(amount *int) {
	if len(os.Args) > 1 {
		newAmount, err := strconv.Atoi(os.Args[1])
		if err != nil {
			logger.LogError("producer — failed to convert argument to string", err)
			*amount = 10
		} else {
			*amount = newAmount
		}
	}
}

func getOrders(amount int) []models.Order {
	var orders []models.Order
	for range amount {
		orders = append(orders, createOrder())
	}
	return orders
}

func createOrder() models.Order {
	var order models.Order

	order.OrderUID = newOderUID()
	order.TrackNumber = newTrackNumber()
	order.Entry = "WBIL"
	order.Delivery = createDelivery()
	order.Payment = createPayment(order)
	order.Items = createItems(order)
	order.Locale = newLocale()
	order.InternalSignature = ""
	order.CustomerID = newCustomerID()
	order.DeliveryService = newDeliveryService()
	order.ShardKey = newShardKey()
	order.SmID = newSmID()
	order.DateCreated = newDateCreated()
	order.OofShard = newOofShard()

	return order
}

func createDelivery() models.Delivery {
	var delivery models.Delivery

	delivery.Name = newName()
	delivery.Phone = newPhone()
	delivery.Zip = newZip()
	delivery.City = newCity()
	delivery.Address = newAddress()
	delivery.Region = newRegion()
	delivery.Email = newEmail()

	return delivery
}

func createPayment(order models.Order) models.Payment {
	var payment models.Payment

	payment.Transaction = order.OrderUID
	payment.RequestID = ""
	payment.Currency = newCurrency()
	payment.Provider = "wbpay"
	payment.Amount = newAmount()
	payment.PaymentDT = 1637907727 // can't be bothered
	payment.Bank = newBank()
	payment.DeliveryCost = newDeliveryCost()
	payment.GoodsTotal = newGoodsTotal()
	payment.CustomFee = newCustomFee()

	return payment
}

func createItems(order models.Order) []models.Item {
	var items []models.Item
	totalItems := totalItems()
	for range totalItems {
		items = append(items, newItem(order))
	}
	return items
}

func newItem(order models.Order) models.Item {
	var item models.Item

	item.ChrtID = newChrtId()
	item.TrackNumber = order.TrackNumber
	item.Price = newPrice()
	item.Rid = newRid()
	item.Name = newItemName()
	item.Sale = newSale()
	item.Size = newSize()
	item.TotalPrice = newTotalPrice(item.Price, item.Sale)
	item.NmID = newNmId()
	item.Brand = newBrand()
	item.Status = newStatus()

	return item
}

func newOderUID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	orderUID := hex.EncodeToString(bytes) + "test"
	return orderUID
}

func newTrackNumber() string {
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	bytes := make([]byte, 9)
	rand.Read(bytes)
	for i := 0; i < 9; i++ {
		bytes[i] = letters[int(bytes[i])%len(letters)]
	}
	return "WBILM" + string(bytes)
}

func newLocale() string {
	locales := []string{"en", "ru", "de", "zh", "fr", "es", "it", "ja"}
	number := big.NewInt(int64(len(locales)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newLocale — failed to create random number", err)
	}
	return locales[idx.Int64()]
}

func newCustomerID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	id := hex.EncodeToString(bytes)
	return id
}

func newDeliveryService() string {
	names := []string{
		"wildberries",
		"meest",
		"boxberry",
		"hermes",
		"ems",
		"pickpoint",
		"kurier",
		"ozon",
	}
	number := big.NewInt(int64(len(names)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newDeliveryService — failed to create random number", err)
	}
	return names[idx.Int64()]
}

func newShardKey() string {
	max := big.NewInt(10)
	key, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newShardKey — failed to create random number", err)
	}
	return fmt.Sprintf("%d", key.Int64())
}

func newSmID() int {
	max := big.NewInt(100)
	id, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newSmID — failed to create random number", err)
	}
	return int(id.Int64())
}

func newDateCreated() time.Time {
	return time.Now()
}

func newOofShard() string {
	max := big.NewInt(100)
	key, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newOofShard — failed to create random number", err)
	}
	return fmt.Sprintf("%d", key.Int64())
}

func newName() string {
	names := []string{
		"Test Testov",
		"Max Payne",
		"John Doe",
		"Nikita Chetverkin",
		"Martin Molin",
		"Charlie White",
		"Jeremy Jahns",
		"Chris Stuckmann",
		"Harrier Du Bois",
		"Manuel Calavera",
		"Demetrian Titus",
		"Jill Valentine",
		"Isaac Clarke",
		"Gordon Freeman",
	}
	number := big.NewInt(int64(len(names)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newName — failed to create random number", err)
	}
	return names[idx.Int64()]
}

func newPhone() string {
	const digits = "0123456789"
	bytes := make([]byte, 10)
	rand.Read(bytes)
	for i := range bytes {
		bytes[i] = digits[int(bytes[i])%len(digits)]
	}
	return "+9" + string(bytes)
}

func newZip() string {
	const digits = "0123456789"
	bytes := make([]byte, 7)
	rand.Read(bytes)
	for i := range bytes {
		bytes[i] = digits[int(bytes[i])%len(digits)]
	}
	return string(bytes)
}

func newCity() string {
	cities := []string{
		"Raccoon City",
		"Vice City",
		"Los Santos",
		"San Fierro",
		"Las Venturas",
		"Silent Hill",
		"Gotham City",
		"Metropolis",
		"Night City",
		"King's Landing",
		"Bikini Bottom",
		"Sin City",
	}

	number := big.NewInt(int64(len(cities)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newCity — failed to create random number", err)
	}
	return cities[idx.Int64()]
}

func newAddress() string {
	addresses := []string{
		"Ploshad Mira 15",
		"Baker Street 221B",
		"Sesame Street 123",
		"Fleet Street 186",
		"Grove Street 4",
		"Elm Street 1428",
		"Wall Street 12",
		"Privet Drive 4",
		"Mulholland Drive 17",
	}
	number := big.NewInt(int64(len(addresses)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newAddress — failed to create random number", err)
	}
	return addresses[idx.Int64()]
}

func newRegion() string {
	regions := []string{
		"California",
		"Bavaria",
		"Catalonia",
		"Quebec",
		"Scotland",
		"Flanders",
		"Lombardy",
		"Kyushu",
		"Queensland",
		"Auckland",
		"Rio de Janeiro",
		"Buenos Aires Province",
		"Gauteng",
	}

	number := big.NewInt(int64(len(regions)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newRegion — failed to create random number", err)
	}
	return regions[idx.Int64()]
}

func newEmail() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	mail := hex.EncodeToString(bytes) + "@gmail.com"
	return mail
}

func newCurrency() string {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CNY", "RUB", "AUD"}
	number := big.NewInt(int64(len(currencies)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newCurrency — failed to create random number", err)
	}
	return currencies[idx.Int64()]
}

func newAmount() float64 {
	max := big.NewInt(100000)
	amount, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newAmount — failed to create random number", err)
	}
	return float64(amount.Int64()) / 100.0
}

func newBank() string {
	banks := []string{"alpha", "sber", "vtb", "gazprombank", "bank of america", "deutsche bank", "chase", "santander"}
	number := big.NewInt(int64(len(banks)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newBank — failed to create random number", err)
	}
	return banks[idx.Int64()]
}

func newDeliveryCost() float64 {
	max := big.NewInt(5000)
	cost, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newDeliveryCost — failed to create random number", err)
	}
	return float64(cost.Int64()) / 100.0
}

func newGoodsTotal() float64 {
	max := big.NewInt(200000)
	total, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newGoodsTotal — failed to create random number", err)
	}
	return float64(total.Int64()) / 100.0
}

func newCustomFee() float64 {
	max := big.NewInt(10000)
	fee, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newCustomFee — failed to create random number", err)
	}
	return float64(fee.Int64()) / 100.0
}

func totalItems() int64 {
	max := big.NewInt(5)
	total, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("totalItems — failed to create random number", err)
	}
	return total.Int64()
}
func newChrtId() int {
	id, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		logger.LogFatal("newChrtId — failed to create random number", err)
	}
	return int(id.Int64())
}

func newPrice() float64 {
	price, err := rand.Int(rand.Reader, big.NewInt(5000))
	if err != nil {
		logger.LogFatal("newPrice — failed to create random number", err)
	}
	return float64(price.Int64())
}

func newRid() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes) + "test"
}

func newItemName() string {
	items := []string{
		"Lightsaber",
		"The One Ring",
		"Wand",
		"Infinity Gauntlet",
		"Mjolnir",
		"Neuralyzer",
		"Hoverboard",
		"Portal Gun",
		"Gravity Gun",
		"Witcher Medallion",
		"BFG 9000",
		"Dragon Glass Dagger",
		"Potion of Healing",
		"Tesseract",
		"Dragon Egg",
		"Witcher Sword Silver",
		"Witcher Sword Steel",
		"Elder Scroll",
		"Death Note",
		"Wizard Staff",
		"Samurai Armor",
		"Holy Grail Cup",
		"Red Pill",
		"Blue Pill",
		"Necronomicon",
		"Sword of Gryffindor",
		"Voldemort's Horcrux",
		"Arc Reactor",
		"Captain Jack's Compass",
		"Golden Ticket",
		"Heisenberg Hat",
		"Vicodin Bottle",
		"Sherlock's Pipe",
		"Walter Sobchak Bowling Ball",
		"Ghostbusters Proton Pack",
		"Wilson Volleyball",
		"Magic 8-Ball",
		"Big Kahuna Burger",
		"Inception Totem Top",
		"Leon's Plant",
		"Rocky Gloves",
		"Chewbacca Plush",
		"Saw Puzzle Box",
		"Jigsaw Mask",
		"Kubrick Monolith Replica",
		"Neo's Sunglasses",
		"John Wick' Pistol",
	}
	number, err := rand.Int(rand.Reader, big.NewInt(int64(len(items))))
	if err != nil {
		logger.LogFatal("newItemName — failed to generate random index", err)
	}
	return items[number.Int64()]
}

func newSale() int {
	sale, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		logger.LogFatal("newSale — failed to generate random number", err)
	}
	return int(sale.Int64())
}

func newSize() string {
	sizes := []string{"XS", "S", "M", "L", "XL"}
	number, err := rand.Int(rand.Reader, big.NewInt(int64(len(sizes))))
	if err != nil {
		logger.LogFatal("newSize — failed to generate random index", err)
	}
	return sizes[number.Int64()]
}

func newTotalPrice(price float64, sale int) float64 {
	return price * (1 - float64(sale)/100)
}

func newNmId() int {
	number, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		logger.LogFatal("newNmId — failed to create random number", err)
	}
	return int(number.Int64())
}

func newBrand() string {
	brands := []string{"Vivienne Sabo", "Maybelline", "L'Oreal", "NYX", "Revlon"} // change
	number, err := rand.Int(rand.Reader, big.NewInt(int64(len(brands))))
	if err != nil {
		logger.LogFatal("newBrand — failed to generate random index", err)
	}
	return brands[number.Int64()]
}

func newStatus() int {
	statuses := []int{100, 200, 202, 300, 400}
	number, err := rand.Int(rand.Reader, big.NewInt(int64(len(statuses))))
	if err != nil {
		logger.LogFatal("newStatus — failed to generate random index", err)
	}
	return int(number.Int64())
}
