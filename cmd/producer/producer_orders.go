package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/logger"
	"github.com/Pur1st2EpicONE/WBTECH-sample-microservice/internal/models"
)

func GetOrders(amount int, logger logger.Logger) []models.Order {
	var orders []models.Order
	for range amount {
		orders = append(orders, CreateOrder(logger))
	}
	return orders
}

func CreateOrder(logger logger.Logger) models.Order {
	var order models.Order

	order.OrderUID = newOderUID()
	order.TrackNumber = newTrackNumber()
	order.Entry = "WBIL"
	order.Delivery = createDelivery(logger)
	order.Payment = createPayment(order, logger)
	order.Items = createItems(order, logger)
	order.Locale = newLocale(logger)
	order.InternalSignature = ""
	order.CustomerID = newCustomerID()
	order.DeliveryService = newDeliveryService(logger)
	order.ShardKey = newShardKey(logger)
	order.SmID = newSmID(logger)
	order.DateCreated = newDateCreated()
	order.OofShard = newOofShard(logger)

	return order
}

func createDelivery(logger logger.Logger) models.Delivery {
	var delivery models.Delivery

	delivery.Name = newName(logger)
	delivery.Phone = newPhone()
	delivery.Zip = newZip()
	delivery.City = newCity(logger)
	delivery.Address = newAddress(logger)
	delivery.Region = newRegion(logger)
	delivery.Email = newEmail()

	return delivery
}

func createPayment(order models.Order, logger logger.Logger) models.Payment {
	var payment models.Payment

	payment.Transaction = order.OrderUID
	payment.RequestID = ""
	payment.Currency = newCurrency(logger)
	payment.Provider = "wbpay"
	payment.Amount = newAmount(logger)
	payment.PaymentDT = 1637907727 // can't be bothered
	payment.Bank = newBank(logger)
	payment.DeliveryCost = newDeliveryCost(logger)
	payment.GoodsTotal = newGoodsTotal(logger)
	payment.CustomFee = newCustomFee(logger)

	return payment
}

func createItems(order models.Order, logger logger.Logger) []models.Item {
	var items []models.Item
	totalItems := totalItems(logger)
	for range totalItems {
		items = append(items, newItem(order, logger))
	}
	return items
}

func newItem(order models.Order, logger logger.Logger) models.Item {
	var item models.Item

	item.ChrtID = newChrtId(logger)
	item.TrackNumber = order.TrackNumber
	item.Price = newPrice(logger)
	item.Rid = newRid()
	item.Name = newItemName(logger)
	item.Sale = newSale(logger)
	item.Size = newSize(logger)
	item.TotalPrice = newTotalPrice(item.Price, item.Sale)
	item.NmID = newNmId(logger)
	item.Brand = newBrand(logger)
	item.Status = newStatus(logger)

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

func newLocale(logger logger.Logger) string {
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

func newDeliveryService(logger logger.Logger) string {
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

func newShardKey(logger logger.Logger) string {
	max := big.NewInt(10)
	key, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newShardKey — failed to create random number", err)
	}
	return fmt.Sprintf("%d", key.Int64())
}

func newSmID(logger logger.Logger) int {
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

func newOofShard(logger logger.Logger) string {
	max := big.NewInt(100)
	key, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newOofShard — failed to create random number", err)
	}
	return fmt.Sprintf("%d", key.Int64())
}

func newName(logger logger.Logger) string {
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

func newCity(logger logger.Logger) string {
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

func newAddress(logger logger.Logger) string {
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

func newRegion(logger logger.Logger) string {
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

func newCurrency(logger logger.Logger) string {
	currencies := []string{"USD", "EUR", "GBP", "JPY", "CNY", "RUB", "AUD"}
	number := big.NewInt(int64(len(currencies)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newCurrency — failed to create random number", err)
	}
	return currencies[idx.Int64()]
}

func newAmount(logger logger.Logger) float64 {
	max := big.NewInt(100000)
	amount, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newAmount — failed to create random number", err)
	}
	return float64(amount.Int64()) / 100.0
}

func newBank(logger logger.Logger) string {
	banks := []string{"alpha", "sber", "vtb", "gazprombank", "bank of america", "deutsche bank", "chase", "santander"}
	number := big.NewInt(int64(len(banks)))
	idx, err := rand.Int(rand.Reader, number)
	if err != nil {
		logger.LogFatal("newBank — failed to create random number", err)
	}
	return banks[idx.Int64()]
}

func newDeliveryCost(logger logger.Logger) float64 {
	max := big.NewInt(5000)
	cost, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newDeliveryCost — failed to create random number", err)
	}
	return float64(cost.Int64()) / 100.0
}

func newGoodsTotal(logger logger.Logger) float64 {
	max := big.NewInt(200000)
	total, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newGoodsTotal — failed to create random number", err)
	}
	return float64(total.Int64()) / 100.0
}

func newCustomFee(logger logger.Logger) float64 {
	max := big.NewInt(10000)
	fee, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("newCustomFee — failed to create random number", err)
	}
	return float64(fee.Int64()) / 100.0
}

func totalItems(logger logger.Logger) int64 {
	max := big.NewInt(5)
	total, err := rand.Int(rand.Reader, max)
	if err != nil {
		logger.LogFatal("totalItems — failed to create random number", err)
	}
	return total.Int64()
}
func newChrtId(logger logger.Logger) int {
	id, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		logger.LogFatal("newChrtId — failed to create random number", err)
	}
	return int(id.Int64())
}

func newPrice(logger logger.Logger) float64 {
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

func newItemName(logger logger.Logger) string {
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

func newSale(logger logger.Logger) int {
	sale, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		logger.LogFatal("newSale — failed to generate random number", err)
	}
	return int(sale.Int64())
}

func newSize(logger logger.Logger) string {
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

func newNmId(logger logger.Logger) int {
	number, err := rand.Int(rand.Reader, big.NewInt(10000000))
	if err != nil {
		logger.LogFatal("newNmId — failed to create random number", err)
	}
	return int(number.Int64())
}

func newBrand(logger logger.Logger) string {
	brands := []string{
		"Oscorp",
		"Wayne Enterprises",
		"Stark Industries",
		"Weyland-Yutani",
		"Vault-Tec",
		"Aperture Science",
		"Black Mesa",
	}
	number, err := rand.Int(rand.Reader, big.NewInt(int64(len(brands))))
	if err != nil {
		logger.LogFatal("newBrand — failed to generate random index", err)
	}
	return brands[number.Int64()]
}

func newStatus(logger logger.Logger) int {
	statuses := []int{100, 200, 202, 300, 400}
	number, err := rand.Int(rand.Reader, big.NewInt(int64(len(statuses))))
	if err != nil {
		logger.LogFatal("newStatus — failed to generate random index", err)
	}
	return int(number.Int64())
}
