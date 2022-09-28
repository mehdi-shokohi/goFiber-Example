package messages

// Global
const (
	InternalServerError = "مشکلی در سرور رخ داده است"
	BanIp               = "شما به دلیل تکرار بیش از حد مسدود شده اید"
	BlockedUser         = "حساب کاربری شما مسدود شده است"
	InactiveUser        = "حساب کاربری شما فعال نمی باشد"
	InvalidInputForm    = "فرم ارسال شده معتبر نمی باشد"
	RequiredKey         = "ارسال کلید جستجو الزامی است"
)

// authManager
const (
	OnlySendAccessToken  = "فقط اکسس توکن ارسال کنید"
	WrongAccessToken     = "اکسس توکن ارسال شده اشتباه است"
	ExpiredAccessToken   = "توکن شما منقضی شده است"
	UserIsNotExist       = "کاربر وجود ندارد"
	RequiredRefreshToken = "لطفا رفرش توکن را وارد کنید"
	WrongRefreshToken    = "رفرش توکن ارسال شده اشتباه است"
	ExpiredRefreshToken  = "رفرش توکن شما منقضی شده است"
	OnlySendRefreshToken = "فقط رفرش توکن ارسال کنید"
	PleaseLogin          = "لطفا وارد شوید"
	UserIsNotAdmin       = "کاربر ادمین نمی باشد"
	SuccessfullLogin     = "شما با موفقیت وارد شدید"
	SuccessfullLogout    = "شما با موفقیت خارج شدید"
)

// adminUser
const (
	RequiredUsername       = "لطفا نام کاربری را وارد کنید"
	RequiredPassword       = "لطفا رمز عبور را وارد کنید"
	InvalidUsername        = "لطفا نام کاربری معتبر وارد کنید"
	InvalidPassword        = "لطفا رمز عبور معتبر وارد کنید"
	InvalidMelliCode       = "کد ملی وارد شده صحیح نمی باشد"
	InvalidEmail           = "لطفا ایمیل معتبر وارد کنید"
	InvalidMobileNumber    = "لطفا شماره موبایل معتبر وارد کنید"
	InvalidFirstName       = "لطفا نام معتبر وارد کنید"
	InvalidLastName        = "لطفا نام خانوادگی معتبر وارد کنید"
	NotFoundUsername       = "نام کاربری وارد شده وجود ندارد"
	WrongPassword          = "رمز عبور صحیح نمی باشد"
	NotFoundAdminUsername  = "ادمینی با این نام کاربری وجود ندارد"
	DuplicateUsername      = "این نام کاربری تکراری می باشد"
	AdminCanNotChangeAdmin = "ادمین نمیتواند اطلاعات ادمین دیگری را تغییر دهد"

	SuccessfullCreatedUser = "کاربر جدید با موفقیت ایجاد شد"
	SuccessfullEditUser    = "اطلاعات کاربر با موفقیت تغییر یافت"
)

// Email Messages.
const (
	EmailError = "خطایی در ارسال ایمیل به وجود آمد"
)

// Items
const (
	ItemIdRequired = "شناسه آیتم الزامی است"
	InvalidItemID  = "شناسه نامعتبر می باشد"
)

// doctors
const (
	JustOneFieldDoctor = "فقط نام یا کدنظام پزشکی را وارد کنید"
)

// prescription
const (
	RequiredPersonIDPrescription        = "شناسه شخص الزامی است"
	RequiredContractIDPrescription      = "شناسه بیمه الزامی است"
	RequiredDoctorIDPrescription        = "شناسه دکتر الزامی است"
	RequiredDoctorCodePrescription      = "کد دکتر الزامی است"
	RequiredAdmissionPrescription       = "تاریخ پذیرش الزامی است"
	RequiredServiceDatePrescription     = "تاریخ ثبت الزامی است"
	RequiredEffectivePrescription       = "شناسه تاثیر گذار بیمه الزامی است"
	RequiredGradeCodePrescription       = "کد رتبه الزامی است"
	RequiredTotalPricePrescription      = "قیمت کل را وارد کنید"
	RequiredBaseInsurePricePrescription = "قیمت پایه بیمه را وارد کنید"

	RequiredItemId         = "شناسه آیتم لازم است"
	RequiredItemCode       = "کد آیتم لازم است"
	RequiredInvoiceId      = "شناسه صورت حساب الزامی است"
	RequiredRequestCount   = "لطفا مقدار آیتم را وارد کنید"
	RequiredRequestPrice   = "لطفا مبلغ آیتم را وارد کنید"
	RequiredInsuranceOrgId = "لطفا کد سازمان بیمه را وارد کنید"
)

// Insurance Org's
const (
	InsuranceOrgIdNotFound  = "کد سازمان بیمه یافت نشد"
	InsuranceOrgIsNotActive = "این سازمان بیمه فعال نمی باشد"
)

// Passwords
const (
	PasswordsNotMatch                  = "رمز عبور قدیمی با جدید یکی نیست"
	PasswordSuccessfullyEdited         = "رمز عبور با موفقیت تغییر یافت"
	OTPCodeForPasswordSendSuccessfully = "کد یکبار مصرف به ایمیل کاربر ارسال شد"
	RequiredVerificationCode           = "کد یکبار مصرف را وارد کنید"
	PleaseSendRequestAgain             = "مشکلی در فرایند ایمیل به وجود آمد , لطفا دوباره درخواست دهید"
	CodeCorrect                        = "کد وارد شده درست است"
	PleaseWaitForEndTimeDuration       = "لطفا حداکثر تا سه دقیقه صبر کنید سپس اقدام کنید"

	EnterWithOutValidation = "لطفا اول درخواست ارسال کد بدهید"

	CodeIncorrect = "کد وارد شده مطابقت ندارد"
)

// Specializes

const (
	SpecializesNameRequired = "نام تخصص الزامی است"
	SpecializesCodeRequired = "کد تخصص الزامی است"
	HAVNTGRANT = "شما مجوز ندارید"
)
