from django.db import models
from django.contrib.auth.models import User
import datetime
# Create your models here.
class Cert(models.Model):
	token = models.CharField(max_length=36)
	user = models.ForeignKey(User, on_delete=models.CASCADE, related_name="cert_user")
	title = models.CharField(max_length=256)
	location = models.CharField(max_length=256)
	org = models.ForeignKey(User, on_delete=models.CASCADE, related_name="cert_org")
	user_agrees = models.CharField(max_length=10)
	org_filled = models.BooleanField(default=False)
	time_till = models.DateField()
	is_verified = models.BooleanField(default=False)
	
	# def __str__(self):
	#     return f"{self.user.username} {self.org.username} {self.token}"

class Resume(models.Model):
	org = models.CharField(max_length=100)
	user = models.CharField(max_length=100)
	tokens = models.CharField(max_length=1000)
	sent = models.BooleanField(default=False)

	def __str__(self):
		return f"{self.user} to {self.org} with {{len(json.loads(self.tokens))}} fields"