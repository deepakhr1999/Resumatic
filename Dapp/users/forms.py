from django import forms
from django.contrib.auth.models import User
# from django.contrib.auth.forms import UserCreationForm
from .models import *
from kyc.models import *
from django.forms.widgets import TextInput
from django.forms import ModelForm
from django.contrib.auth.forms import UserCreationForm
COMPANY_CHOICES= [
    ('GitHub', 'GitHub'),
    ('Oracle', 'Oracle'),
    ('J P Morgan', 'J P Morgan'),
    ]

LOCATION_CHOICES=[
    ('San Francisco', 'San Francisco'),
    ('Bangalore', 'Bangalore'),
]

# from .models import UserProfile
class UserRegisterForm(UserCreationForm):
    first_name = forms.CharField(widget=forms.TextInput())
    last_name = forms.CharField(widget=forms.TextInput())
    email = forms.EmailField()
    is_org = forms.BooleanField(required=False)
    class Meta:
        model = User
        fields = ['username', 'is_org', 'first_name', 'last_name', 'email', 'password1', 'password2']

class UserUpdateForm(forms.ModelForm):
    email = forms.EmailField()

    class Meta:
        model = User
        fields = ['username', 'email']


class ProfileUpdateForm(forms.ModelForm):
    class Meta:
        model = Profile
        fields = '__all__'

class RegisterCertForm(forms.ModelForm):
    class Meta:
        model = Cert
        fields = ['org', 'title', 'time_till', 'location']

class ProvideCertForm(forms.ModelForm):
    class Meta:
        model = Cert
        fields = ['user', 'title', 'time_till', 'location']

    

