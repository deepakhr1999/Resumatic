from django.db import models
from django.contrib.auth.models import User
import datetime

class Profile(models.Model):
	user = models.OneToOneField(User, on_delete=models.CASCADE)
	dob  = models.DateField(default=datetime.date(2000,1,1))
	is_org = models.BooleanField(default=False)
	image = models.ImageField(default='default.png', upload_to='profile_pics')
	is_verified = models.BooleanField(default=False)
	def __str__(self):
		return f'{self.user.username} Profile'

	def save(self, *args, **kwargs):
		super().save()


